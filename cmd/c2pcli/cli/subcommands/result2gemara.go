package subcommands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer2"
	"github.com/revanite-io/sci/layer4"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func NewResult2Gemara(logger hclog.Logger) *cobra.Command {
	option := NewOptions()
	command := &cobra.Command{
		Use:   "result2gemara",
		Short: "Transform policy result artifacts to Gemara Layer 4 Evaluations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := option.Complete(cmd, logger); err != nil {
				return err
			}
			return runResult2Gemara(cmd.Context(), option)
		},
	}

	fs := command.Flags()
	fs.StringP("out", "o", "-", "path to output directory. Use '-' for stdout. Default '-'.")
	BindGemaraFlags(fs)
	BindPluginFlags(fs)
	return command
}

func runResult2Gemara(ctx context.Context, option *Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}
	manager, err := framework.NewPluginManager(frameworkConfig)
	if err != nil {
		return err
	}

	// We need to load evaluation plans
	config, err := getPolicy(option.Policy)
	if err != nil {
		return err
	}

	for i := range config.Catalogs {
		// Lazily load evals
		config.Catalogs[i].Loader = func() (*layer2.Layer2, error) {
			var l2Catalog layer2.Layer2
			data, err := os.ReadFile(config.Catalogs[i].CatalogID)
			if err != nil {
				return nil, err
			}

			err = yaml.Unmarshal(data, &l2Catalog)
			if err != nil {
				return nil, err
			}
			return &l2Catalog, nil
		}
	}

	inputContext, err := actions.NewContextFromCatalogRefs(config.Catalogs...)
	if err != nil {
		return err
	}

	foundPlugins, err := manager.FindRequestedPlugins(inputContext.RequestedProviders())
	if err != nil {
		return err
	}

	var configSelections framework.PluginConfig = func(pluginID plugin.ID) map[string]string {
		return option.Plugins[pluginID.String()]
	}
	launchedPlugins, err := manager.LaunchPolicyPlugins(foundPlugins, configSelections)
	// Defer clean before returning an error to avoid unterminated processes
	defer manager.Clean()
	if err != nil {
		return err
	}

	filePathResults := filepath.Join(option.EvalDir, "results")
	filePathResults = filepath.Clean(filePathResults)
	err = os.MkdirAll(filePathResults, 0750)
	if err != nil {
		return err
	}

	results, err := actions.AggregateResults(ctx, inputContext, launchedPlugins)
	if err != nil {
		return err
	}

	plans, err := actions.Layer4FromResults(ctx, inputContext, option.Name, results)
	if err != nil {
		return err
	}

	for _, ref := range plans {

		data, err := MarshalConfigWithFunction(ref.Plan)
		if err != nil {
			return err
		}

		// Write the resulting evaluation for each service to a new file
		filePath := filepath.Clean(filepath.Join(filePathResults, fmt.Sprintf("%s.yml", ref.Service)))
		if err := os.WriteFile(filePath, data, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

type Marshalable struct {
	layer4.Layer4
}

// Temporary until unmarshalling is done upstream

func (mc Marshalable) MarshalYAML() (interface{}, error) {
	outputMap := make(map[string]interface{})
	outputMap["catalog_id"] = mc.CatalogID
	outputMap["start_time"] = mc.StartTime
	outputMap["end_time"] = mc.EndTime
	outputMap["corrupted_state"] = mc.CorruptedState

	controlEvals := []map[string]interface{}{}
	for _, controlEval := range mc.ControlEvaluations {
		evalMap := make(map[string]interface{})
		evalMap["control_id"] = controlEval.ControlID
		assessments := []map[string]interface{}{}
		for _, assessment := range controlEval.Assessments {
			assessmentMap := make(map[string]interface{})
			assessmentMap["requirement_id"] = assessment.RequirementID
			methods := []map[string]interface{}{}
			for _, method := range assessment.Methods {
				methodMap := make(map[string]interface{})
				methodMap["name"] = method.Name
				methodMap["description"] = method.Description
				methodMap["run"] = method.Run
				if method.Result != nil {
					methodMap["result"] = map[string]interface{}{
						"status": method.Result.Status,
					}
				}
				methods = append(methods, methodMap)
			}
			assessmentMap["methods"] = methods
			assessments = append(assessments, assessmentMap)
		}
		evalMap["assessments"] = assessments
		controlEvals = append(controlEvals, evalMap)
	}
	outputMap["evaluations"] = controlEvals

	// Return the map, which yaml.Marshal will then convert into YAML.
	return outputMap, nil
}

// MarshalConfigWithFunction takes a Config struct and returns its YAML representation
// by leveraging the custom MarshalYAML implementation.
func MarshalConfigWithFunction(eval *layer4.Layer4) ([]byte, error) {
	marshalable := Marshalable{*eval}

	yamlData, err := yaml.Marshal(marshalable)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config to YAML using custom marshaler: %w", err)
	}
	return yamlData, nil
}

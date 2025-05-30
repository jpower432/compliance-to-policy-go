/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package eval

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer4"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/options"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func NewResult2SCI(logger hclog.Logger) *cobra.Command {
	option := options.NewOptions()

	command := &cobra.Command{
		Use:   "result2sci",
		Short: "Transform policy result artifacts to SCI Layer 4 Evaluations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := option.Complete(cmd, logger); err != nil {
				return err
			}
			return runResult2SCI(cmd.Context(), option)
		},
	}

	fs := command.Flags()
	fs.StringP("out", "o", "-", "path to output directory. Use '-' for stdout. Default '-'.")
	options.BindSCIFlags(fs)
	options.BindPluginFlags(fs)
	return command
}

func runResult2SCI(ctx context.Context, option *options.Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}
	manager, err := framework.NewPluginManager(frameworkConfig)
	if err != nil {
		return err
	}

	getPolicy, err := GetPolicy(option.Policy)
	if err != nil {
		return err
	}

	err = findRefs(getPolicy, option.EvalDir)
	if err != nil {
		return err
	}

	inputContext, err := actions.NewContextFromRefs(getPolicy.Refs...)
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

	for _, ref := range getPolicy.Refs {
		provider := launchedPlugins[ref.PluginID]
		if err := actions.Evaluate(ctx, inputContext, &ref, provider); err != nil {
			return err
		}

		data, err := MarshalConfigWithFunction(ref.Plan)
		if err != nil {
			return err
		}

		// Write the resulting evaluation for each service to a new file
		filePath := filepath.Clean(filepath.Join(filePathResults, fmt.Sprintf("%s.yml", ref.Service)))
		return os.WriteFile(filePath, data, os.ModePerm)
	}
	return nil
}

func findRefs(getPolicy policy.Policy, evalDir string) error {
	for _, catalogPath := range getPolicy.Catalogs {
		catalog, err := getCatalog(catalogPath)
		if err != nil {
			return err
		}
		// Set loaders
		for i := range getPolicy.Refs {
			// Plans are under the plugin name <plugin-id>-<catalog-id>.yml
			filePath := filepath.Clean(filepath.Join(evalDir, fmt.Sprintf("%s-%s.yml", getPolicy.Refs[i].PluginID, catalog.Metadata.Id)))
			getPolicy.Refs[i].Loader = func() (*layer4.Layer4, error) {
				var l4Eval layer4.Layer4
				file, err := os.Open(filePath)
				if err != nil {
					return nil, err
				}
				decoder := yaml.NewDecoder(file)
				err = decoder.Decode(&l4Eval)
				if err != nil {
					return nil, err
				}
				return &l4Eval, nil
			}
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

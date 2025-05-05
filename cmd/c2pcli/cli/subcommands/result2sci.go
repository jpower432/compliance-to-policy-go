package subcommands

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func NewResult2SCI(logger hclog.Logger) *cobra.Command {
	options := NewOptions()
	options.logger = logger

	command := &cobra.Command{
		Use:   "result2sci",
		Short: "Transform policy result artifacts to SCI Layer 4 Evaluations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(cmd); err != nil {
				return err
			}
			if err := options.Validate(); err != nil {
				return err
			}
			return runResult2SCI(cmd.Context(), options)
		},
	}

	fs := command.Flags()
	fs.String(Catalog, "", "Path to Layer 2 SCI Catalog")
	BindPluginFlags(fs)

	return command
}

func runResult2SCI(ctx context.Context, option *Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	plan, _, err := createOrGetPlan(ctx, option)
	if err != nil {
		return err
	}
	inputContext, err := Context(plan)
	if err != nil {
		return err
	}

	manager, err := framework.NewPluginManager(frameworkConfig)
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

	results, err := actions.AggregateResults(ctx, inputContext, launchedPlugins)
	if err != nil {
		return err
	}

	var controls []layer2.Control
	catalog, err := getCatalog(option.Catalog)
	if err != nil {
		return err
	}
	for _, family := range catalog.ControlFamilies {
		controls = append(controls, family.Controls...)
	}

	evals, err := actions.Evaluate(ctx, inputContext, controls, results)
	if err != nil {
		return err
	}

	for _, eval := range evals {
		data, err := yaml.Marshal(eval)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(data))
	}
	return nil
}

func getCatalog(filepath string) (layer2.Catalog, error) {
	var catalog layer2.Catalog
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return catalog, err
	}
	err = yaml.Unmarshal(yamlFile, &catalog)
	if err != nil {
		return catalog, err
	}
	return catalog, nil
}

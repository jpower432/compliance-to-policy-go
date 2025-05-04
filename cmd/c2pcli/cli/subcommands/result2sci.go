package subcommands

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

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

	evals, err := actions.Evaluate(ctx, inputContext, results)
	if err != nil {
		return err
	}

	for _, eval := range evals {
		fmt.Println(eval)
	}
	return nil
}

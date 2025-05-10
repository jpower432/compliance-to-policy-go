package subcommands

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func NewSCI2Policy(logger hclog.Logger) *cobra.Command {
	options := NewOptions()
	options.logger = logger

	command := &cobra.Command{
		Use:   "oscal2policy",
		Short: "Transform OSCAL to policy artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(cmd); err != nil {
				return err
			}
			if err := options.Validate(); err != nil {
				return err
			}
			return runSCI2Policy(cmd.Context(), options)
		},
	}
	fs := command.Flags()
	fs.String(Catalog, "", "Path to Layer 2 SCI Catalog")
	BindPluginFlags(fs)
	return command
}

// Intended output - code generation - plans and policy as code for the plan
func runSCI2Policy(ctx context.Context, option *Options) error {
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

	err = actions.GeneratePolicy(ctx, inputContext, launchedPlugins)
	if err != nil {
		return err
	}

	return nil
}

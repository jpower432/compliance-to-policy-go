package subcommands

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/pipelines"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation/gemara"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation/oscal"
)

var generateTemplatesFlag bool

func NewCompliance2Policy(logger hclog.Logger) *cobra.Command {
	options := NewOptions()

	command := &cobra.Command{
		Use:   "compliance2policy",
		Short: "Transform compliance artifacts to policy artifacts.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(cmd, logger); err != nil {
				return err
			}
			return run2Policy(cmd.Context(), options, args[0])
		},
	}

	fs := command.Flags()
	BindPluginFlags(fs)
	BindGemaraFlags(fs)
	BindOSCALFlags(fs)
	fs.BoolVar(&generateTemplatesFlag, "generate-templates", false, "Set to true to generate evaluation plan templates")
	return command
}

func run2Policy(ctx context.Context, option *Options, kind string) error {
	var provider evaluation.Provider
	switch kind {
	case "gemara":
		if generateTemplatesFlag {
			err := gemara.GenerateTemplates(option.Policy, option.EvalDir)
			if err != nil {
				return err
			}
			fmt.Printf("Fill out created templtes in %s\n", option.EvalDir)
			return nil
		}
		gemaraProvider, err := gemara.NewGemaraValidatorFromFile(option.Policy, option.EvalDir)
		if err != nil {
			return err
		}
		provider = gemaraProvider

	case "oscal":
		if generateTemplatesFlag {
			err := oscal.GenerateTemplates(ctx, option.logger, option.OSCALOptions.Definition, option.Plan, option.Name)
			if err != nil {
				return err
			}
			return nil
		}
		oscalProvider, err := oscal.NewOSCALValidationFromFile(option.Plan, option.logger)
		if err != nil {
			return err
		}
		provider = oscalProvider
	}

	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	var configSelections framework.PluginConfig = func(pluginID plugin.ID) map[string]string {
		return option.Plugins[pluginID.String()]
	}
	pipeline := pipelines.Pipeline{Planner: provider}
	return pipeline.GenerateWorkflow(ctx, frameworkConfig, configSelections)
}

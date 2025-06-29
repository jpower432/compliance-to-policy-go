package subcommands

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/pipelines"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation/gemara"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation/oscal"
)

func NewResult2Compliance(logger hclog.Logger) *cobra.Command {
	option := NewOptions()
	command := &cobra.Command{
		Use:   "result2compliance",
		Short: "Transform policy result artifacts to a compliance artifact.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := option.Complete(cmd, logger); err != nil {
				return err
			}
			return run2Results(cmd.Context(), option, args[0])
		},
	}

	fs := command.Flags()
	fs.StringP("out", "o", "-", "path to output directory. Use '-' for stdout. Default '-'.")
	BindGemaraFlags(fs)
	BindOSCALFlags(fs)
	BindPluginFlags(fs)
	return command
}

func run2Results(ctx context.Context, option *Options, kind string) error {
	var provider evaluation.Provider
	var output string
	switch kind {
	case "gemara":
		gemaraProvider, err := gemara.NewGemaraValidatorFromFile(option.Policy, option.EvalDir)
		if err != nil {
			return err
		}
		provider = gemaraProvider
		resultsPath := filepath.Join(option.EvalDir, "results")
		cleanedPath := filepath.Clean(resultsPath)
		err = os.MkdirAll(cleanedPath, 0700)
		if err != nil {
			return err
		}
		output = cleanedPath

	case "oscal":
		oscalProvider, err := oscal.NewOSCALValidationFromFile(option.Plan, option.logger)
		if err != nil {
			return err
		}
		provider = oscalProvider
		output = option.Output
	}

	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	var configSelections framework.PluginConfig = func(pluginID plugin.ID) map[string]string {
		return option.Plugins[pluginID.String()]
	}
	pipeline := pipelines.Pipeline{Planner: provider, ReportPath: output}
	return pipeline.ResultsPipeline(ctx, frameworkConfig, configSelections)
}

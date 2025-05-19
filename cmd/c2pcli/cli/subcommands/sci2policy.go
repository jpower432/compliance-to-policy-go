/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package subcommands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer4"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func NewSCI2Policy(logger hclog.Logger) *cobra.Command {
	options := NewOptions()
	options.logger = logger

	command := &cobra.Command{
		Use:   "sci2policy",
		Short: "Transform SCI to policy artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(cmd); err != nil {
				return err
			}
			if err := options.Validate(); err != nil {
				return err
			}
			return runSCI2Policy(options)
		},
	}
	fs := command.Flags()
	// Replace with Layer 3 policy
	fs.StringP("plugin-dir", "p", "c2p-plugins", "path to plugin directory. Defaults to `c2p-plugins`.")
	fs.StringP(AssessmentPlan, "a", "", "path to assessment-plan.json. This option cannot be used with --component-definition.")
	fs.StringP("out", "o", "-", "path to output directory. Use '-' for stdout. Default '-'.")
	fs.StringP(ConfigPath, "c", "c2p-config.yaml", "path to the configuration for the C2P CLI.")
	return command
}

// Intended output - code generation - plans and policy as code for the plan
func runSCI2Policy(option *Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	policy, err := getPolicy(option.Plan)
	if err != nil {
		return err
	}

	// Set loaders
	// TODO: Find a way to make this optional
	for i := range policy.Refs {
		// Lazily load evals
		policy.Refs[i].Loader = func() (*layer4.Layer4, error) {
			var l4Eval layer4.Layer4
			return &l4Eval, nil
		}
	}

	inputContext, err := actions.NewContextFromRefs(policy.Refs...)
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

	for _, cat := range policy.Catalogs {
		catalog, err := getCatalog(cat)
		if err != nil {
			return err
		}
		for _, ref := range policy.Refs {
			eval, err := actions.GenerateEvaluation(catalog, launchedPlugins[ref.PluginID])
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(eval)
			if err != nil {
				return err
			}

			out := option.Output
			if out == "-" {
				fmt.Fprintln(os.Stdout, string(data))
			} else {
				err := os.MkdirAll(out, os.ModePerm)
				if err != nil {
					return err
				}
				filePath := filepath.Clean(filepath.Join(out, fmt.Sprintf("%s.yml", ref.Service)))
				return os.WriteFile(filePath, data, os.ModePerm)
			}
		}
	}

	return nil
}

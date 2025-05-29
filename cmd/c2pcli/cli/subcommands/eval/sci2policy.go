/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package eval

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer4"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/options"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func NewSCI2Policy(logger hclog.Logger) *cobra.Command {
	option := options.NewOptions()

	command := &cobra.Command{
		Use:   "sci2policy",
		Short: "Transform SCI to policy artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := option.Complete(cmd, logger); err != nil {
				return err
			}
			return runSCI2Policy(option)
		},
	}
	fs := command.Flags()
	options.BindPluginFlags(fs)
	options.BindSCIFlags(fs)
	return command
}

// Intended output - code generation - plans and policy as code for the plan
func runSCI2Policy(option *options.Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	policy, err := GetPolicy(option.Policy)
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

			err = os.MkdirAll(option.EvalDir, os.ModePerm)
			if err != nil {
				return err
			}
			filePath := filepath.Clean(filepath.Join(option.EvalDir, fmt.Sprintf("%s.yml", ref.Service)))
			return os.WriteFile(filePath, data, os.ModePerm)
		}
	}

	return nil
}

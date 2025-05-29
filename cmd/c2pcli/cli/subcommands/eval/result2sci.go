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

	policy, err := GetPolicy(option.Policy)
	if err != nil {
		return err
	}

	// Set loaders
	for i := range policy.Refs {
		// Lazily load evals
		policy.Refs[i].Loader = func() (*layer4.Layer4, error) {
			var l4Eval layer4.Layer4
			filePath := filepath.Clean(filepath.Join(option.EvalDir, fmt.Sprintf("%s.yml", policy.Refs[i].Service)))
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

	inputContext, err := actions.NewContextFromRefs(policy.Refs...)
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

	for _, ref := range policy.Refs {
		provider := launchedPlugins[ref.PluginID]
		if err := actions.Evaluate(ctx, inputContext, &ref, provider); err != nil {
			return err
		}

		data, err := yaml.Marshal(ref.Plan)
		if err != nil {
			return err
		}

		err = os.MkdirAll(option.EvalDir, os.ModeDir)
		if err != nil {
			return err
		}
		filePath := filepath.Clean(filepath.Join(option.EvalDir, fmt.Sprintf("%s.yml", ref.Service)))
		return os.WriteFile(filePath, data, os.ModePerm)
	}
	return nil
}

/*
Copyright 2023 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package subcommands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
)

func NewOSCAL2Policy() *cobra.Command {
	opts := NewOptions()

	command := &cobra.Command{
		Use:   "oscal2policy",
		Short: "Transform OSCAL to policy artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runOSCAL2Policy(cmd.Context(), opts)
		},
	}

	opts.AddFlags(command.Flags())

	return command
}

func runOSCAL2Policy(ctx context.Context, options *Options) error {
	var pluginsPath *string
	if options.PluginsPath != "" {
		pluginsPath = &options.PluginsPath
	}
	frameworkConfig, err := Config(options.ComponentDefinition, pluginsPath)
	if err != nil {
		return err
	}

	settings, err := Settings(options, frameworkConfig)
	if err != nil {
		return err
	}

	manager, err := framework.NewPluginManager(frameworkConfig)
	if err != nil {
		return err
	}
	foundPlugins, err := manager.FindRequestedPlugins()
	if err != nil {
		return err
	}

	var configSelections map[string]map[string]string
	launchedPlugins, err := manager.LaunchPolicyPlugins(foundPlugins, configSelections)
	if err != nil {
		return err
	}

	err = manager.GeneratePolicy(ctx, launchedPlugins, settings.AllSettings())
	if err != nil {
		return err
	}

	return nil
}

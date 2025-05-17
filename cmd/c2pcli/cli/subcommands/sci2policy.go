/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package subcommands

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

var pluginName plugin.ID

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
			return runSCI2Policy(options)
		},
	}
	fs := command.Flags()
	fs.String(Catalog, "", "Path to Layer 2 SCI Catalog")
	fs.StringVar((*string)(&pluginName), "plugin", "", "Plugin to use")
	BindPluginFlags(fs)
	return command
}

// Intended output - code generation - plans and policy as code for the plan
func runSCI2Policy(option *Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	catalog, err := getCatalog(option.Catalog)
	if err != nil {
		return err
	}

	manager, err := framework.NewPluginManager(frameworkConfig)
	if err != nil {
		return err
	}
	foundPlugins, err := manager.FindRequestedPlugins([]plugin.ID{pluginName})
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

	eval, err := actions.GenerateEvaluation(catalog, launchedPlugins[pluginName])
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(eval)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(os.Stdout, data)
	if err != nil {
		return err
	}

	return nil
}

func getCatalog(filepath string) (layer2.Layer2, error) {
	var catalog layer2.Layer2
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

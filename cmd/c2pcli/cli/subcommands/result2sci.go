/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package subcommands

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

type Policy struct {
	refs []actions.PlanRef `yaml:"refs"`
}

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
	// Replace with Layer 3 policy
	fs.String(AssessmentPlan, "", "Path to L3 policy")
	BindPluginFlags(fs)

	return command
}

func runResult2SCI(ctx context.Context, option *Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}
	manager, err := framework.NewPluginManager(frameworkConfig)
	if err != nil {
		return err
	}

	policy, err := getPolicy(option.Plan)
	if err != nil {
		return err
	}

	inputContext, err := actions.NewContextFromRefs(policy.refs...)
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

	for _, ref := range policy.refs {
		provider := launchedPlugins[ref.PluginID]
		rs, err := actions.Evaluate(ctx, inputContext, ref, provider)
		if err != nil {
			return err
		}
		data, err := yaml.Marshal(ref.Plan)
		if err != nil {
			return err
		}
		fmt.Println(rs.ID)
		fmt.Fprintln(os.Stdout, string(data))
	}

	return nil
}

// TODO: Also load plans here
func getPolicy(filepath string) (Policy, error) {
	var policy Policy
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return policy, err
	}
	err = yaml.Unmarshal(yamlFile, &policy)
	if err != nil {
		return policy, err
	}
	return policy, nil
}

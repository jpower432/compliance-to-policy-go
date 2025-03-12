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

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
)

func NewResult2OSCAL() *cobra.Command {
	opts := NewResultOptions(NewOptions())

	command := &cobra.Command{
		Use:   "result2oscal",
		Short: "Transform policy result artifact to OSCAL Assessment Results.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runResult2Policy(cmd.Context(), opts)
		},
	}

	opts.AddFlags(command.Flags())

	return command
}

func runResult2Policy(ctx context.Context, options *ResultOptions) error {
	var pluginsPath *string
	if options.PluginsPath != "" {
		pluginsPath = &options.PluginsPath
	}
	frameworkConfig, err := Config(options.ComponentDefinition, pluginsPath)
	if err != nil {
		return err
	}

	settings, err := Settings(options.Options, frameworkConfig)
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

	var configSelection map[string]map[string]string
	launchedPlugins, err := manager.LaunchPolicyPlugins(foundPlugins, configSelection)
	if err != nil {
		return err
	}

	results, err := manager.AggregateResults(ctx, launchedPlugins, settings.AllSettings())
	if err != nil {
		return err
	}

	reporter, err := framework.NewReporter(frameworkConfig)
	if err != nil {
		return err
	}

	assessmentResults, err := reporter.GenerateAssessmentResults(ctx, "REPLACE_ME", settings, results)
	oscalModels := oscalTypes.OscalModels{
		AssessmentResults: &assessmentResults,
	}

	err = pkg.WriteObjToJsonFile(options.OutputPath, oscalModels)
	if err != nil {
		return err
	}
	return nil
}

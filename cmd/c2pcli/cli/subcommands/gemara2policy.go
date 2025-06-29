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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer2"
	"github.com/revanite-io/sci/layer4"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var generateTemplatesFlag bool

func NewGemara2Policy(logger hclog.Logger) *cobra.Command {
	options := NewOptions()

	command := &cobra.Command{
		Use:   "gemara2policy",
		Short: "Transform Gemara Layer 3 policy to policy artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(cmd, logger); err != nil {
				return err
			}
			return runGemara2Policy(cmd.Context(), options)
		},
	}

	fs := command.Flags()
	BindPluginFlags(fs)
	BindGemaraFlags(fs)
	fs.BoolVar(&generateTemplatesFlag, "generate-templates", false, "Set to true to generate evaluation plan templates")
	return command
}

func runGemara2Policy(ctx context.Context, option *Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	getPlans, err := getPolicy(option.Policy)
	if err != nil {
		return err
	}

	if generateTemplatesFlag {
		_, err := generateTemplates(&getPlans, option.EvalDir)
		if err != nil {
			return err
		}
		fmt.Printf("Fill out created templtes in %s\n", option.EvalDir)
		return nil
	}

	for i := range getPlans.Catalogs {
		// Lazily load evals
		getPlans.Catalogs[i].Loader = func() (*layer2.Layer2, error) {
			catalog, err := getCatalog(getPlans.Catalogs[i].CatalogID)
			return &catalog, err
		}
	}

	inputContext, err := actions.NewContextFromCatalogRefs(getPlans.Catalogs...)
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

func generateTemplates(getPolicy *policy.Config, evalDir string) (foundRefs []policy.PlanRef, err error) {
	// FIXME: Duplicates need to be removed one per validator and Layer 2 catalog
	for i := range getPolicy.Catalogs {
		catalog, err := getCatalog(getPolicy.Catalogs[i].CatalogID)
		if err != nil {
			return foundRefs, err
		}
		getPolicy.Catalogs[i].Catalog = &catalog

		// Set loaders
		for _, ref := range getPolicy.Catalogs[i].Plans {
			// Config are under the plugin name <plugin-id>-<catalog-id>.yml
			filePath := filepath.Clean(filepath.Join(evalDir, fmt.Sprintf("%s-%s.yml", ref.PluginID, catalog.Metadata.Id)))
			if _, err := os.Stat(filePath); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					if err := generateNewEvalTemplate(catalog, filePath); err != nil {
						return foundRefs, err
					}
					// Do not generate policy if creating a template from scratch
					continue
				}
				return foundRefs, err
			}
			foundRefs = append(foundRefs, ref)
		}
	}
	return foundRefs, nil
}

func generateNewEvalTemplate(catalog layer2.Layer2, filePath string) error {
	eval := layer4.NewEvaluation(catalog)
	data, err := yaml.Marshal(eval)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, os.ModePerm)
}

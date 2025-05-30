/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package eval

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/revanite-io/sci/layer2"
	"github.com/revanite-io/sci/layer4"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/options"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/policy"
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
			return runSCI2Policy(cmd.Context(), option)
		},
	}
	fs := command.Flags()
	options.BindPluginFlags(fs)
	options.BindSCIFlags(fs)
	return command
}

// Intended output - code generation - plans and policy as code for the plan
func runSCI2Policy(ctx context.Context, option *options.Options) error {
	frameworkConfig, err := Config(option)
	if err != nil {
		return err
	}

	getPolicy, err := GetPolicy(option.Policy)
	if err != nil {
		return err
	}

	err = os.MkdirAll(option.EvalDir, os.ModeDir)
	if err != nil {
		return err
	}

	foundRefs, err := generateTemplates(getPolicy, option.EvalDir)
	if err != nil {
		return err
	}

	if foundRefs == nil || len(foundRefs) == 0 {
		fmt.Printf("No evaluation plans found. Fill out created templtes in %s\n", option.EvalDir)
		return nil
	}

	inputContext, err := actions.NewContextFromRefs(foundRefs...)
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

func generateTemplates(getPolicy policy.Policy, evalDir string) (foundRefs []policy.PlanRef, err error) {
	// FIXME: Duplicates need to be removed one per validator and Layer 2 catalog
	for _, catalogPath := range getPolicy.Catalogs {
		catalog, err := getCatalog(catalogPath)
		if err != nil {
			return foundRefs, err
		}
		// Set loaders
		for _, ref := range getPolicy.Refs {
			// Plans are under the plugin name <plugin-id>-<catalog-id>.yml
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
			ref.Loader = func() (*layer4.Layer4, error) {
				var l4Eval layer4.Layer4
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

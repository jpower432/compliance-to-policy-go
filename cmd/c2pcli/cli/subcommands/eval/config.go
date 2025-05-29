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

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/oscal-compass/oscal-sdk-go/transformers"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/revanite-io/sci/layer2"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/options"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/policy"
)

// Config returns a populated C2PConfig for the CLI to use.
func Config(option *options.Options) (*framework.C2PConfig, error) {
	c2pConfig := framework.DefaultConfig()
	pluginsPath := option.PluginDir
	if pluginsPath != "" {
		c2pConfig.PluginDir = pluginsPath
		c2pConfig.PluginManifestDir = pluginsPath
	}
	// Set logger
	c2pConfig.Logger = option.Logger()
	return c2pConfig, nil
}

func Context(ap *oscalTypes.AssessmentPlan) (*actions.InputContext, error) {
	if ap.LocalDefinitions == nil || ap.LocalDefinitions.Activities == nil || ap.AssessmentAssets.Components == nil {
		return nil, errors.New("error converting component definition to assessment plan")
	}

	var allComponents []components.Component
	for _, component := range *ap.AssessmentAssets.Components {
		compAdapter := components.NewSystemComponentAdapter(component)
		allComponents = append(allComponents, compAdapter)
	}

	inputCtx, err := actions.NewContext(allComponents)
	if err != nil {
		return nil, err
	}

	apSettings := settings.NewAssessmentActivitiesSettings(*ap.LocalDefinitions.Activities)
	inputCtx.Settings = apSettings

	return inputCtx, nil
}

// createOrGetPlan will load an OSCAL Assessment Plan if detected from the options for return the loaded plan and file location.
// If no plan is detected, it is created from an OSCAL Component Definition for a given framework name.
func createOrGetPlan(ctx context.Context, option *options.Options) (*oscalTypes.AssessmentPlan, string, error) {
	if option.Plan != "" {
		plan, err := loadPlan(option.Plan)
		if err != nil {
			return nil, "", fmt.Errorf("error loading assessment plan: %w", err)
		}
		return plan, option.Plan, nil
	}
	compDef, err := loadCompDef(option.Definition)
	if err != nil {
		return nil, "", fmt.Errorf("error loading component definition: %w", err)
	}

	ap, err := transformers.ComponentDefinitionsToAssessmentPlan(ctx, []oscalTypes.ComponentDefinition{compDef}, option.Name)
	if err != nil {
		return nil, "", err
	}

	return ap, "REPLACE_ME", nil
}

func loadCompDef(path string) (oscalTypes.ComponentDefinition, error) {
	file, err := os.Open(path)
	if err != nil {
		return oscalTypes.ComponentDefinition{}, err
	}
	defer file.Close()
	compDef, err := models.NewComponentDefinition(file, validation.NewSchemaValidator())
	if err != nil {
		return oscalTypes.ComponentDefinition{}, err
	}

	if compDef == nil {
		return oscalTypes.ComponentDefinition{}, errors.New("component definition cannot be nil")
	}
	return *compDef, nil
}

func loadPlan(path string) (*oscalTypes.AssessmentPlan, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	plan, err := models.NewAssessmentPlan(file, validation.NewSchemaValidator())
	if err != nil {
		return nil, err
	}
	return plan, nil
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

// TODO: Also load plans here

func GetPolicy(filepath string) (policy.Policy, error) {
	var p policy.Policy
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

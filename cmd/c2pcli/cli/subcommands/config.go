/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package subcommands

import (
	"errors"
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/revanite-io/sci/layer2"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// Config returns a populated C2PConfig for the CLI to use.
func Config(option *Options) (*framework.C2PConfig, error) {
	c2pConfig := framework.DefaultConfig()
	pluginsPath := option.PluginDir
	if pluginsPath != "" {
		c2pConfig.PluginDir = pluginsPath
		c2pConfig.PluginManifestDir = pluginsPath
	}
	// Set logger
	c2pConfig.Logger = option.logger
	return c2pConfig, nil
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

func getPolicy(filepath string) (policy.Config, error) {
	var p policy.Config
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

/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package subcommands

import (
	"errors"
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/oscal-compass/oscal-sdk-go/validation"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/config"
)

// Config returns a populated C2PConfig for the CLI to use.
func Config(option *Options) (*config.C2PConfig, error) {
	c2pConfig := config.DefaultConfig()
	pluginsPath := option.PluginDir
	if pluginsPath != "" {
		c2pConfig.PluginDir = pluginsPath
	}
	// Set logger
	c2pConfig.Logger = option.logger
	return c2pConfig, nil
}

func Target(option *Options) (*config.Target, error) {
	compDef, err := loadCompDef(option.Definition)
	if err != nil {
		return nil, err
	}
	return config.NewTargetFromComponentDefinition(compDef)
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

// Settings returns extracted compliance settings from a given component definition implementation using the C2PConfig.
func Settings(option *Options) (*settings.ImplementationSettings, error) {
	var implementation []oscalTypes.ControlImplementationSet
	compDef, err := loadCompDef(option.Definition)
	if err != nil {
		return nil, err
	}
	for _, cp := range *compDef.Components {
		if cp.ControlImplementations != nil {
			implementation = append(implementation, *cp.ControlImplementations...)
		}
	}
	return settings.Framework(option.Name, implementation)
}

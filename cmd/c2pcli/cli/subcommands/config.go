package subcommands

import (
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/config"
)

func Config(componentPath string, pluginsPaths *string) (*config.C2PConfig, error) {
	c2pConfig := config.DefaultConfig()
	if pluginsPaths != nil {
		c2pConfig.PluginDir = *pluginsPaths
	}

	compDef, err := loadCompDef(componentPath)
	if err != nil {
		return nil, err
	}
	c2pConfig.ComponentDefinitions = []oscalTypes.ComponentDefinition{*compDef}
	return c2pConfig, nil
}

func loadCompDef(path string) (*oscalTypes.ComponentDefinition, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	compDef, err := generators.NewComponentDefinition(file)
	if err != nil {
		return nil, err
	}
	return compDef, nil
}

func Settings(options *Options, frameworkConfig *config.C2PConfig) (*settings.ImplementationSettings, error) {
	var implementation []oscalTypes.ControlImplementationSet
	for _, comp := range frameworkConfig.ComponentDefinitions {
		for _, cp := range *comp.Components {
			if cp.ControlImplementations != nil {
				implementation = append(implementation, *cp.ControlImplementations...)
			}
		}
	}
	return settings.Framework(options.Name, implementation)
}

/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package framework

import (
	"fmt"

	"github.com/hashicorp/go-hclog"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/action"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/config"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// PluginManager manages the plugin lifecycle and compliance-to-policy
// workflows.
type PluginManager struct {
	// pluginDir is the location to search for plugins.
	pluginDir string
	// clientFactory is the function used to
	// create new plugin clients.
	clientFactory plugin.ClientFactoryFunc
	// logger for the PluginManager
	log hclog.Logger
}

// NewPluginManager creates a new instance of a PluginManager from a C2PConfig that can be used to
// interact with supported plugins.
//
// It supports the plugin lifecycle with the following methods:
//   - Finding and initializing plugins: FindRequestedPlugins() and LaunchPolicyPlugins()
//   - Clean/Stop - Clean()
func NewPluginManager(cfg *config.C2PConfig) (*PluginManager, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &PluginManager{
		pluginDir:     cfg.PluginDir,
		clientFactory: plugin.ClientFactory(cfg.Logger),
		log:           cfg.Logger,
	}, nil
}

// FindRequestedPlugins retrieves information for the plugins that have been requested
// in the C2PConfig and returns the plugin manifests for use with LaunchPolicyPlugins().
func (m *PluginManager) FindRequestedPlugins(target *action.Target, pluginType string) (plugin.Manifests, error) {
	m.log.Info(fmt.Sprintf("Searching for plugins in %s", m.pluginDir))

	pluginManifests, err := plugin.FindPlugins(
		m.pluginDir,
		plugin.WithProviderIds(target.RequiredPlugins()),
		plugin.WithPluginType(pluginType),
	)
	if err != nil {
		return pluginManifests, err
	}
	m.log.Debug(fmt.Sprintf("Found %d matching plugins", len(pluginManifests)))
	return pluginManifests, nil
}

// LaunchPolicyPlugins launches requested plugins and configures each plugin to make it ready for use with defined plugin workflows.
// The plugin is configured based on default options and given options.
// Given options are represented by config.PluginConfig.
func (m *PluginManager) LaunchPolicyPlugins(manifests plugin.Manifests, pluginConfig config.PluginConfig) (map[string]policy.Aggregator, error) {
	pluginsByIds := make(map[string]policy.Aggregator)
	for _, manifest := range manifests {
		policyPlugin, err := plugin.NewPolicyPlugin(manifest, m.clientFactory)
		if err != nil {
			return pluginsByIds, err
		}
		pluginsByIds[manifest.ID] = policyPlugin
		m.log.Debug(fmt.Sprintf("Launched plugin %s", manifest.ID))
		m.log.Debug(fmt.Sprintf("Gathering configuration options for %s", manifest.ID))

		// Get all the base configuration
		if len(manifest.Configuration) > 0 {
			if err := m.configurePlugin(policyPlugin, manifest, pluginConfig); err != nil {
				return pluginsByIds, fmt.Errorf("failed to configure plugin %s: %w", manifest.ID, err)
			}
		}
	}
	return pluginsByIds, nil
}

// LaunchGeneratorPlugins launches requested plugins and configures each plugin to make it ready for use with defined plugin workflows.
// The plugin is configured based on default options and given options.
// Given options are represented by config.PluginConfig.
func (m *PluginManager) LaunchGeneratorPlugins(manifests plugin.Manifests, pluginConfig config.PluginConfig) (map[string]policy.Generator, error) {
	pluginsByIds := make(map[string]policy.Generator)
	for _, manifest := range manifests {
		policyPlugin, err := plugin.NewGeneratorPlugin(manifest, m.clientFactory)
		if err != nil {
			return pluginsByIds, err
		}
		pluginsByIds[manifest.ID] = policyPlugin
		m.log.Debug(fmt.Sprintf("Launched plugin %s", manifest.ID))
		m.log.Debug(fmt.Sprintf("Gathering configuration options for %s", manifest.ID))

		// Get all the base configuration
		if len(manifest.Configuration) > 0 {
			if err := m.configurePlugin(policyPlugin, manifest, pluginConfig); err != nil {
				return pluginsByIds, fmt.Errorf("failed to configure plugin %s: %w", manifest.ID, err)
			}
		}
	}
	return pluginsByIds, nil
}

func (m *PluginManager) configurePlugin(policyPlugin policy.Provider, manifest plugin.Manifest, pluginConfig config.PluginConfig) error {
	selections := pluginConfig(manifest.ID)
	if selections == nil {
		selections = make(map[string]string)
		m.log.Debug("No overrides set for plugin %s, using defaults...", manifest.ID)
	}
	configMap, err := manifest.ResolveOptions(selections)
	if err != nil {
		return err
	}
	if err := policyPlugin.Configure(configMap); err != nil {
		return err
	}
	return nil
}

// Clean deletes managed instances of plugin clients that have been created using LaunchPolicyPlugins.
// This will remove all clients launched with the plugin.ClientFactoryFunc.
func (m *PluginManager) Clean() {
	m.log.Debug("Cleaning launched plugins")
	plugin.Cleanup()
}

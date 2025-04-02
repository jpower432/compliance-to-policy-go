/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// ServeConfig defines the configuration for plugin
// registration.
type ServeConfig struct {
	PluginSet map[string]plugin.Plugin
	Logger    hclog.Logger
}

// Register a set of implemented plugins.
// This function should be called last during plugin initialization in the main function.
func Register(config ServeConfig) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         config.PluginSet,
		Logger:          config.Logger,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

// Cleanup clean up all plugin clients created by the ClientFactory.
var Cleanup func() = plugin.CleanupClients

// ClientFactoryFunc defines a function signature for creating
// new go-plugin clients.
type ClientFactoryFunc func(manifest Manifest) (*plugin.Client, error)

// ClientFactory returns a factory function for creating new plugin-specific
// clients with consistent plugin config settings.
//
// The returned factory function takes a Manifest object as input and returns
// a new plugin client configured with the specified logger, allowed protocols,
// and security settings.
func ClientFactory(logger hclog.Logger) ClientFactoryFunc {
	return func(manifest Manifest) (*plugin.Client, error) {
		manifestSum, err := hex.DecodeString(manifest.Checksum)
		if err != nil {
			return nil, err
		}
		config := &plugin.ClientConfig{
			HandshakeConfig: Handshake,
			Logger:          logger.Named(manifest.ID),
			// Enabling this will ensure that client.Kill() is run when this is cleaned up.
			Managed:          true,
			AutoMTLS:         true,
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			Cmd:              exec.Command(manifest.ExecutablePath),
			Plugins:          SupportedPlugins,
			SecureConfig: &plugin.SecureConfig{
				Checksum: manifestSum,
				Hash:     sha256.New(),
			},
		}

		client := plugin.NewClient(config)
		return client, nil
	}
}

// NewPolicyPlugin dispenses a new instance of a policy plugin.
func NewPolicyPlugin(pluginManifest Manifest, createClient ClientFactoryFunc) (policy.Aggregator, error) {
	client, err := createClient(pluginManifest)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin client for %s: %w", pluginManifest.ID, err)
	}
	rpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin client for %s: %w", pluginManifest.ID, err)
	}

	raw, err := rpcClient.Dispense(PVPPluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to dispense plugin %s: %w", pluginManifest.ID, err)
	}

	p := raw.(policy.Aggregator)
	return p, nil
}

// NewGeneratorPlugin dispenses a new instance of a policy plugin.
func NewGeneratorPlugin(pluginManifest Manifest, createClient ClientFactoryFunc) (policy.Generator, error) {
	client, err := createClient(pluginManifest)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin client for %s: %w", pluginManifest.ID, err)
	}
	rpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin client for %s: %w", pluginManifest.ID, err)
	}

	raw, err := rpcClient.Dispense(GenerationPluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to dispense plugin %s: %w", pluginManifest.ID, err)
	}

	p := raw.(policy.Generator)
	return p, nil
}

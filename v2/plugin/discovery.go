package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

const pluginDir = "./plugin-dir"

// Selector stores plugin.Config types by the plugin ID.
type Selector map[string]Config

// FindPlugins searches for plugin locations by provider id while filtering on plugin type.
func (s Selector) FindPlugins(providerIds []string, pluginType string) (pluginsPath map[string]string, err error) {
	ex, err := os.Executable()
	if err != nil {
		return pluginsPath, err
	}
	pluginRootPath := filepath.Join(filepath.Dir(ex), pluginDir)
	pluginsPath = map[string]string{}

	for _, providerId := range providerIds {
		pluginConfig, ok := s[providerId]
		if !ok {
			return pluginsPath, fmt.Errorf("no plugin config found for provider %s", providerId)
		}
		for _, implementedTypes := range pluginConfig.ImplementedTypes {
			if implementedTypes == pluginType {
				pluginPath := fmt.Sprintf("%s/%s", pluginRootPath, providerId)
				pluginsPath[providerId] = pluginPath
				break
			}
		}
	}
	return pluginsPath, nil
}

// Config store information about installed plugins.
type Config struct {
	Version          string   `json:"version"`
	ImplementedTypes []string `json:"type"`
}

// NewPolicyClient dispenses a new instance of a policy plugin.
func NewPolicyClient(pluginBinaryPath string) (providers.PolicyProvider, *hplugin.Client, error) {
	client := createClient(pluginBinaryPath)
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, err
	}

	raw, err := rpcClient.Dispense(PVPPluginName)
	if err != nil {
		return nil, nil, err
	}

	// Pass back the client so that the caller can kill the child process
	p := raw.(providers.PolicyProvider)
	return p, client, nil
}

// NewRemediationClient dispenses a new instance of a remediation plugin.
func NewRemediationClient(pluginBinaryPath string) (providers.RemediationProvider, *hplugin.Client, error) {
	client := createClient(pluginBinaryPath)
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, err
	}

	raw, err := rpcClient.Dispense(RemediationPluginName)
	if err != nil {
		return nil, nil, err
	}

	// Pass back the client so that the caller can kill the child process
	p := raw.(providers.RemediationProvider)
	return p, client, nil
}

// NewGenerationClient dispenses a new instance of a generation plugin.
func NewGenerationClient(pluginBinaryPath string) (providers.GenerationProvider, *hplugin.Client, error) {
	client := createClient(pluginBinaryPath)
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, err
	}

	raw, err := rpcClient.Dispense(GenerationPluginName)
	if err != nil {
		return nil, nil, err
	}

	// Pass back the client so that the caller can kill the child process
	p := raw.(providers.GenerationProvider)
	return p, client, nil
}

// createClient creates a single instance of a plugin client
func createClient(pluginBinaryPath string) *hplugin.Client {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	cmd := exec.Command(pluginBinaryPath)
	client := hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig:  Handshake,
		Plugins:          SupportedPlugins,
		Cmd:              cmd,
		Logger:           logger,
		SyncStdout:       os.Stdout,
		SyncStderr:       os.Stderr,
		AllowedProtocols: []hplugin.Protocol{hplugin.ProtocolGRPC},
	})

	return client
}

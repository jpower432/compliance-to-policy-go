/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
)

const (
	// Only validation components are plugins
	pluginComponentType = "validation"
	// DefaultPluginPath default location c2p will look for plugins
	DefaultPluginPath = "c2p-plugins"
)

// C2PConfig represents configuration options for the C2P framework.PluginManager.
type C2PConfig struct {
	// PluginDir is the directory where the PluginManager searches
	// for installed plugins.
	PluginDir string
	// Logger is the logging implementation used in the PluginManager and
	// plugin clients.
	Logger hclog.Logger
}

var defaultLogger = hclog.New(&hclog.LoggerOptions{
	Name:   "c2p",
	Output: os.Stdout,
	Level:  hclog.Info,
})

// DefaultConfig returns the default configuration.
func DefaultConfig() *C2PConfig {
	return &C2PConfig{
		PluginDir: DefaultPluginPath,
		Logger:    defaultLogger,
	}
}

// Validate returns an error if C2PConfig has invalid fields.
func (c *C2PConfig) Validate() error {
	// Sanitize the plugin directory input
	c.PluginDir = strings.TrimSpace(c.PluginDir)
	c.PluginDir = filepath.Clean(c.PluginDir)
	if strings.TrimSpace(c.PluginDir) == "" {
		return fmt.Errorf("plugin directory cannot be empty")
	}
	if _, err := os.Stat(c.PluginDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("plugin directory %s does not exist: %w", c.PluginDir, err)
		}
		return err
	}
	if c.Logger == nil {
		c.Logger = defaultLogger
	}
	return nil
}

// PluginConfig is a function signature that returns configuration
// option key, value pairs for a given plugin id.
type PluginConfig func(string) map[string]string

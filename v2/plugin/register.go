package plugin

import hplugin "github.com/hashicorp/go-plugin"

// Register a new plugin or set of plugins
func Register(plugins map[string]hplugin.Plugin) {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         plugins,
		GRPCServer:      hplugin.DefaultGRPCServer,
	})
}

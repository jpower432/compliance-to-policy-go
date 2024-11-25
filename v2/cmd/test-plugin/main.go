package main

import (
	"encoding/json"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

var _ providers.PolicyProvider = (*MyExamplePlugin)(nil)
var _ providers.GenerationProvider = (*MyExamplePlugin)(nil)

type MyExamplePlugin struct{}

func (p MyExamplePlugin) GetSchema() ([]byte, error) {
	return nil, nil
}

func (p MyExamplePlugin) UpdateConfiguration(message json.RawMessage) error {
	fmt.Println("I have been configured")
	return nil
}

func (p MyExamplePlugin) Generate(rules providers.Policy) error {
	fmt.Println("I have been generated")
	return nil
}

func (p MyExamplePlugin) GetResults() (providers.PVPResult, error) {
	fmt.Println("I have been scanned")
	return providers.PVPResult{
		ObservationsByCheck: []providers.ObservationByCheck{
			{
				Title:       "example",
				Description: "example",
				Methods:     []string{"AUTOMATED"},
				CheckID:     "etcd_peer_key_file",
			},
		},
	}, nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins: map[string]hplugin.Plugin{
			plugin.PVPPluginName:        &plugin.PVPPlugin{Impl: MyExamplePlugin{}},
			plugin.GenerationPluginName: &plugin.GeneratorPlugin{Impl: MyExamplePlugin{}},
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}

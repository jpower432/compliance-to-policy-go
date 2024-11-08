package main

import (
	"encoding/json"
	"fmt"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/observations"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var _ policy.Provider = &MyExamplePlugin{}

type MyExamplePlugin struct{}

func (p MyExamplePlugin) GetSchema() ([]byte, error) {
	return nil, nil
}

func (p MyExamplePlugin) UpdateConfiguration(message json.RawMessage) error {
	fmt.Println("I have been configured")
	return nil
}

func (p MyExamplePlugin) Generate(rules plan.Policy) error {
	fmt.Println("I have been generated")
	return nil
}

func (p MyExamplePlugin) GetResults() (observations.PVPResult, error) {
	fmt.Println("I have been scanned")
	return observations.PVPResult{
		ObservationsByCheck: []*observations.ObservationByCheck{
			{
				Title:       "example",
				Description: "example",
				Methods:     []string{"AUTOMATED"},
				CheckID:     "me",
			},
		},
	}, nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins: map[string]hplugin.Plugin{
			plugin.PVPPluginName: &plugin.PVPPlugin{Impl: MyExamplePlugin{}},
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}

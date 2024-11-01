package main

import (
	"fmt"
	"github.com/hashicorp/go-plugin"

	shared "github.com/oscal-compass/compliance-to-policy/v2/plugin"
)

type MyExamplePlugin struct{}

func (p MyExamplePlugin) GetSchema() error {
	fmt.Println("I have given a schema")
	return nil
}

func (p MyExamplePlugin) UpdateConfiguration() error {
	fmt.Println("I have been configured")
	return nil
}

func (p MyExamplePlugin) Generate() error {
	fmt.Println("I have generated")
	return nil
}

func (p MyExamplePlugin) GetResults() error {
	return fmt.Errorf("I have results")
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"pvp": &shared.PVPPlugin{Impl: MyExamplePlugin{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

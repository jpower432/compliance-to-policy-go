package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"

	shared "github.com/oscal-compass/compliance-to-policy/v2/plugin"
)

func run() error {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("PVP_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("pvp")
	if err != nil {
		return err
	}

	pvp := raw.(shared.PolicyEngine)
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "generate":
		err := pvp.Generate()
		if err != nil {
			return err
		}

	case "scan":
		err := pvp.GetResults()
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("'scan' and 'generate' are valid, given: %q", os.Args[0])
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

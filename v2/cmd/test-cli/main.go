package main

import (
	"fmt"

	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

func run() error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	client := hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins:         plugin.Map,
		Cmd:             exec.Command("sh", "-c", os.Getenv("PVP_PLUGIN")),
		Logger:          logger,
		AllowedProtocols: []hplugin.Protocol{
			hplugin.ProtocolNetRPC, hplugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(plugin.PVPPluginName)
	if err != nil {
		return err
	}

	pvp := raw.(policy.Engine)
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "generate":
		err := pvp.Generate(policy.Policy{
			Rules: []policy.RuleObject{
				policy.RuleObject{
					RuleId:               "some-id",
					RuleDescription:      "some description",
					PolicyId:             "some-id",
					ParameterId:          "some-id",
					ParameterDescription: "some description",
				},
			},
		})
		if err != nil {
			return err
		}

	case "scan":
		results, err := pvp.GetResults()
		if err != nil {
			return err
		}
		for _, result := range results.Observations {
			fmt.Println(result.Description)
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

package main

import (
	"fmt"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/rules"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"

	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
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

	pvp := raw.(policy.Provider)
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "generate":
		err := pvp.Generate(rules.Policy{
			RuleSets: []*rules.RuleSet{
				&rules.RuleSet{
					RuleID:          "some-id",
					RuleDescription: "some description",
					CheckID:         "some-id",
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
		for _, result := range results.ObservationsByCheck {
			fmt.Printf(fmt.Sprintf("Result %s processed\n", result.Title))
			for _, s := range result.Subjects {
				fmt.Printf("Subject %s result\n", s.Result.String())
			}
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

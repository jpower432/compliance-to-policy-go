package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/report"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/rules"
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
		reporter := report.New(plan.Decomposer{})
		assessmentResult, err := reporter.ToOSCAL(results)
		if err != nil {
			return err
		}

		var b bytes.Buffer
		jsonEncoder := json.NewEncoder(&b)
		jsonEncoder.SetIndent("", "  ")

		if err := jsonEncoder.Encode(assessmentResult); err != nil {
			return err
		}

		if err := os.WriteFile("./assessment-results.json", b.Bytes(), 0600); err != nil {
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

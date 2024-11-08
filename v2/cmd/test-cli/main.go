package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan/rulefinder"
	"os"
	"os/exec"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/report"
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
		SyncStdout:      os.Stdout,
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

	compDefPath := os.Getenv("COMPDEF_PATH")
	compDefJSON, err := os.ReadFile(compDefPath)
	if err != nil {
		return err
	}

	var oscalModels oscalTypes_1_1_2.OscalModels
	dec := json.NewDecoder(bytes.NewBuffer(compDefJSON))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&oscalModels); err != nil {
		return err
	}

	if oscalModels.ComponentDefinition == nil {
		return fmt.Errorf("No component definition found in %s\n", os.Getenv("COMPDEF_PATH"))
	}

	ruleFinder, err := rulefinder.NewMemoryWithCD(*oscalModels.ComponentDefinition)
	if err != nil {
		return err
	}
	newPlan := plan.NewPlan(compDefPath, ruleFinder)

	pvp := raw.(policy.Provider)
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "generate":

		newPolicy, err := newPlan.GetPolicy()
		if err != nil {
			return err
		}

		if err := pvp.Generate(newPolicy); err != nil {
			return err
		}

	case "scan":
		results, err := pvp.GetResults()
		if err != nil {
			return err
		}
		reporter := report.New(newPlan)
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

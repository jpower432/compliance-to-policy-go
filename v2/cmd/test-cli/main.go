package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/rules"

	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func run() error {

	ctx := context.Background()
	compDefPath := os.Getenv("COMPDEF_PATH")
	compDefFile, err := os.Open(compDefPath)
	if err != nil {
		return err
	}

	definition, err := generators.NewComponentDefinition(compDefFile)
	if err != nil {
		return err
	}

	if definition.Components == nil {
		return fmt.Errorf("no component definition found")
	}

	ruleFinder, err := rules.NewMemoryStoreWithComponents(*definition.Components)
	if err != nil {
		return err
	}
	newPlan := oscal.NewPlan(compDefPath, ruleFinder)

	// Resolve all the validation component information
	titleByIds := make(map[string]string)
	providerPolicies := make(map[string]oscal.Policy)
	providerIds := make([]string, 0)
	for _, component := range *definition.Components {
		if component.Type == "validation" {
			id := strings.ToLower(component.Title)
			titleByIds[id] = component.Title
			providerIds = append(providerIds, id)
			providerPolicy, err := newPlan.GetPolicyForComponent(ctx, component.Title)
			if err != nil {
				return err
			}
			providerPolicies[id] = providerPolicy
		}
	}

	// Hard code map of plugins and capabilities
	pluginSelector := plugin.Selector{
		"example": {
			ImplementedTypes: []string{
				plugin.GenerationPluginName,
				plugin.PVPPluginName,
			},
		},
	}

	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "generate":
		plugins, err := pluginSelector.FindPlugins(providerIds, plugin.GenerationPluginName)
		if err != nil {
			return err
		}

		clients := make([]*hplugin.Client, 0, len(plugins))
		for providerId, pluginPath := range plugins {
			generator, client, err := plugin.NewGenerationClient(pluginPath)
			if err != nil {
				return err
			}
			clients = append(clients, client)

			// get the provider ids here to grab the policy
			componentTitle := titleByIds[providerId]
			policy := providerPolicies[componentTitle]

			if err := generator.Generate(policy); err != nil {
				return err
			}
		}

		// Kill child processes
		for _, client := range clients {
			client.Kill()
		}

	case "scan":
		plugins, err := pluginSelector.FindPlugins(providerIds, plugin.GenerationPluginName)
		if err != nil {
			return err
		}

		clients := make([]*hplugin.Client, 0, len(plugins))
		allResults := make([]oscal.PVPResult, 0, len(plugins))
		for _, pluginPath := range plugins {
			pvp, client, err := plugin.NewPolicyClient(pluginPath)
			if err != nil {
				return err
			}
			clients = append(clients, client)
			results, err := pvp.GetResults()
			if err != nil {
				return err
			}
			allResults = append(allResults, results)
		}

		reporter := oscal.NewReporter(newPlan)
		assessmentResult, err := reporter.ToOSCAL(ctx, allResults)
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

		// Kill child processes
		for _, client := range clients {
			client.Kill()
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

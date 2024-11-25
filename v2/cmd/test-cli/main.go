package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/rules"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
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

	ruleFinder, err := rules.NewMemoryStoreFromComponents(*definition.Components)
	if err != nil {
		return err
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

	pluginsManager := framework.NewManager(ruleFinder, pluginSelector)

	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "generate":
		err := pluginsManager.Index(ctx, *definition)
		if err != nil {
			return err
		}

		err = pluginsManager.TransformToPolicy(ctx)
		if err != nil {
			return err
		}

		pluginsManager.Stop()

	case "scan":
		err := pluginsManager.Index(ctx, *definition)
		if err != nil {
			return err
		}

		assessmentResult, err := pluginsManager.AggregateResults(ctx)
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

		pluginsManager.Stop()

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

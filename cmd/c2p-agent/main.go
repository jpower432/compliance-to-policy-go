/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2p-agent/agentkit"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var planConfig, archivistaURL string
	flag.StringVar(&archivistaURL, "archvista-url", "localhost:8080", "URL for Archivista")
	flag.StringVar(&planConfig, "plan", "", "Location for plan config")
	flag.Parse()

	c2pConfig := framework.DefaultConfig()
	manager, err := framework.NewPluginManager(c2pConfig)
	if err != nil {
		return err
	}

	var planRef actions.PlanRef
	file, err := os.Open(planConfig)
	if err != nil {
		return err
	}
	planDecoder := yaml.NewDecoder(file)
	err = planDecoder.Decode(&planRef)
	if err != nil {
		return err
	}
	// TODO: Load eval from disk for now. Eventually it could be imported.
	// TODO: Plugin configuration will need to loading so the plugin knows where
	// generated policy is located if not centrally deployed.

	mn, err := manager.FindRequestedPlugins([]plugin.ID{planRef.PluginID})
	if err != nil {
		return err
	}

	providers, err := manager.LaunchPolicyPlugins(mn, nil)
	defer manager.Clean()
	if err != nil {
		return err
	}

	provider := providers[planRef.PluginID]
	agent := agentkit.NewAgent(provider, planRef)
	err = agent.Run(ctx, agentkit.RunWithExporterURL(archivistaURL))
	if err != nil {
		return err
	}
	return nil
}

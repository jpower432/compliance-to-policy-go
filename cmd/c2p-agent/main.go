/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2p-agent/agentkit"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pluginName := os.Getenv("C2P_PLUGIN")
	c2pConfig := framework.DefaultConfig()
	manager, err := framework.NewPluginManager(c2pConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mn, err := manager.FindRequestedPlugins([]plugin.ID{plugin.ID(pluginName)})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	providers, err := manager.LaunchPolicyPlugins(mn, nil)
	defer manager.Clean()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var provider policy.Provider
	for _, p := range providers {
		provider = p
		break
	}

	url := os.Getenv("ARCHIVISTA_URL")

	// TODO: Load from config
	plan := actions.PlanRef{}
	agent := agentkit.NewAgent(provider, plan)
	err = agent.Run(ctx, agentkit.RunWithExporterURL(url))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package framework

import (
	"context"
	"strings"

	oscaltypes112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

// TODO[jpower432]: Look at implementation concurrent operations here

// Manager executes plugin operations based on OSCAL input.
// Also, this should be functionally equivalent to Python C2P class
type Manager struct {
	store    rules.Store
	selector plugin.Selector

	clients          map[string]*hplugin.Client
	titleByIds       map[string]string
	providerPolicies map[string]providers.Policy
	providerIds      []string
}

// NewManager creates a new framework.Manager
func NewManager(store rules.Store, selector plugin.Selector) *Manager {
	return &Manager{
		store:            store,
		selector:         selector,
		clients:          make(map[string]*hplugin.Client),
		titleByIds:       make(map[string]string),
		providerPolicies: make(map[string]providers.Policy),
		providerIds:      []string{},
	}
}

// Index adds new OSCAL-based configuration to the Manager.
func (m *Manager) Index(ctx context.Context, componentDefinition oscaltypes112.ComponentDefinition) error {
	// Resolve all the validation component information
	for _, component := range *componentDefinition.Components {
		if component.Type == "validation" {
			id := strings.ToLower(component.Title)
			m.titleByIds[id] = component.Title
			m.providerIds = append(m.providerIds, id)
			providerPolicy, err := m.getPolicyForComponent(ctx, component.Title)
			if err != nil {
				return err
			}
			m.providerPolicies[id] = providerPolicy
		}
	}
	return nil
}

func (m *Manager) AggregateResults(ctx context.Context) (assessmentResults oscaltypes112.AssessmentResults, err error) {
	plugins, err := m.selector.FindPlugins(m.providerIds, plugin.PVPPluginName)
	if err != nil {
		return assessmentResults, err
	}

	allResults := make([]providers.PVPResult, 0, len(plugins))
	for providerId, pluginPath := range plugins {
		pvp, client, err := plugin.NewPolicyClient(pluginPath)
		if err != nil {
			return assessmentResults, err
		}
		m.clients[providerId] = client

		results, err := pvp.GetResults()
		if err != nil {
			return assessmentResults, err
		}
		allResults = append(allResults, results)
	}

	reporter := NewReporter(m.store)

	return reporter.ToOSCAL(ctx, "./assessment-plan.json", allResults)
}

func (m *Manager) TransformToPolicy(ctx context.Context) error {
	plugins, err := m.selector.FindPlugins(m.providerIds, plugin.GenerationPluginName)
	if err != nil {
		return err
	}

	for providerId, pluginPath := range plugins {
		generator, client, err := plugin.NewGenerationClient(pluginPath)
		if err != nil {
			return err
		}
		m.clients[providerId] = client

		// get the provider ids here to grab the policy
		componentTitle := m.titleByIds[providerId]
		policy := m.providerPolicies[componentTitle]

		if err := generator.Generate(policy); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Stop() {
	// Kill child processes
	for id, client := range m.clients {
		client.Kill()
		delete(m.clients, id)
	}
}

func (m *Manager) getPolicyForComponent(ctx context.Context, componentTitle string) (providers.Policy, error) {
	collectedRules, err := m.store.FindByComponent(ctx, componentTitle)
	if err != nil {
		return providers.Policy{}, err
	}
	// Change for parameter slice
	parameters := make([]extensions.Parameter, 0, len(collectedRules))
	for _, rule := range collectedRules {
		if rule.Rule.Parameter == nil {
			continue
		}
		parameters = append(parameters, *rule.Rule.Parameter)
	}
	policy := providers.Policy{RuleSets: collectedRules, Parameters: parameters}
	return policy, nil
}

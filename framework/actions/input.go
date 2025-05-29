/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

const pluginComponentType = "validation"

var ErrMissingProvider = errors.New("missing title for provider")

// InputContext is used to configure action behavior from parsed OSCAL documents.
type InputContext struct {
	// requestedProviders stores resolved provider IDs to the original component title from
	// parsed OSCAL Components.
	//
	// The original component title is needed to get information for the validation
	// component in the rules.Store (which provides input for the corresponding policy.Provider
	// id).
	requestedProviders map[plugin.ID]string
	// rulesStore contains indexed information about parsed RuleSets
	// which can be searched for the corresponding component title.
	rulesStore rules.Store
	// Settings define adjustable rule settings parsed from framework-specific implementation
	Settings settings.Settings
}

// NewContext returns an InputContext for the given OSCAL Components.
func NewContext(components []components.Component) (*InputContext, error) {
	inputCtx := &InputContext{
		requestedProviders: make(map[plugin.ID]string),
	}
	for _, comp := range components {
		if comp.Type() == pluginComponentType {
			pluginId, err := GetPluginIDFromComponent(comp)
			if err != nil {
				return inputCtx, err
			}
			inputCtx.requestedProviders[pluginId] = comp.Title()
		}
	}
	store, err := DefaultStore(components)
	if err != nil {
		return inputCtx, err
	}
	inputCtx.rulesStore = store
	return inputCtx, nil
}

// NewContextFromRefs returns an InputContext for the given Layer 3 Policy.
func NewContextFromRefs(refs ...policy.PlanRef) (*InputContext, error) {
	inputCtx := &InputContext{
		requestedProviders: make(map[plugin.ID]string),
	}
	for _, ref := range refs {
		inputCtx.requestedProviders[ref.PluginID] = string(ref.PluginID)
	}
	store, err := newRefStore(refs...)
	if err != nil {
		return inputCtx, err
	}
	inputCtx.rulesStore = store
	inputCtx.Settings = store.Settings()
	return inputCtx, nil
}

// RequestedProviders returns the provider ids requested in the parsed input.
func (t *InputContext) RequestedProviders() []plugin.ID {
	var requestedIds []plugin.ID
	for id := range t.requestedProviders {
		requestedIds = append(requestedIds, id)
	}
	return requestedIds
}

// ProviderTitle return the original component Title for the given provider id.
func (t *InputContext) ProviderTitle(providerId plugin.ID) (string, error) {
	title, ok := t.requestedProviders[providerId]
	if !ok {
		return "", fmt.Errorf("%s:%w", providerId, ErrMissingProvider)
	}
	return title, nil
}

// Store returns the underlying rules.Store with indexed RuleSets.
func (t *InputContext) Store() rules.Store {
	return t.rulesStore
}

// DefaultStore returns a default rules.MemoryStore with indexed information from the given
// Components.
func DefaultStore(allComponents []components.Component) (*rules.MemoryStore, error) {
	store := rules.NewMemoryStore()
	err := store.IndexAll(allComponents)
	if err != nil {
		return store, err
	}
	return store, nil
}

// GetPluginIDFromComponent returns the normalized plugin identifier defined by the OSCAL Component
// of type "validation".
func GetPluginIDFromComponent(component components.Component) (plugin.ID, error) {
	title := strings.TrimSpace(component.Title())
	if title == "" {
		return "", fmt.Errorf("component is missing a title")
	}

	title = strings.ToLower(title)
	id := plugin.ID(title)
	if !id.Validate() {
		return "", fmt.Errorf("invalid plugin id %s", title)
	}
	return id, nil
}

var _ rules.Store = (*refStore)(nil)

type refStore struct {
	// nodes saves the rule ID map keys, which are used with
	// the other fields.
	nodes map[string]extensions.RuleSet
	// ByCheck store a mapping between the checkId and its parent
	// ruleId
	byCheck map[string]string

	// Below contains maps that store information by component and
	// component types to form RuleSet with the correct context.

	// rulesByComponent stores the component title of any component
	// mapped to any relevant rules.
	rulesByComponent map[string]map[string]struct{}
}

func newRefStore(refs ...policy.PlanRef) (*refStore, error) {
	store := &refStore{
		nodes:            make(map[string]extensions.RuleSet),
		byCheck:          make(map[string]string),
		rulesByComponent: make(map[string]map[string]struct{}),
	}

	for _, ref := range refs {
		if ref.Plan == nil {
			err := ref.Load()
			if err != nil {
				return store, err
			}
		}
		if err := store.indexRef(ref); err != nil {
			return store, err
		}
	}
	return store, nil
}

func (p *refStore) indexRef(ref policy.PlanRef) error {
	ruleIds := make(map[string]struct{})
	for _, controlEvals := range ref.Plan.ControlEvaluations {
		for _, as := range controlEvals.Assessments {
			ruleSet := extensions.RuleSet{
				Rule: extensions.Rule{
					ID:          as.RequirementID,
					Description: as.RequirementID,
				},
			}
			for _, method := range as.Methods {
				p.byCheck[method.Name] = as.RequirementID
				ruleSet.Checks = append(ruleSet.Checks, extensions.Check{
					ID:          method.Name,
					Description: method.Description,
				})
			}
			p.nodes[as.RequirementID] = ruleSet
			ruleIds[as.RequirementID] = struct{}{}
		}
	}
	p.rulesByComponent[ref.PluginID.String()] = ruleIds
	p.rulesByComponent[ref.Service] = ruleIds
	return nil
}

func (p *refStore) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	ruleSet, ok := p.nodes[ruleID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("rule %q: %w", ruleID, rules.ErrRuleNotFound)
	}
	return ruleSet, nil
}

func (p *refStore) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	ruleId, ok := p.byCheck[checkID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("failed to find rule for check %q: %w", checkID, rules.ErrRuleNotFound)
	}
	return p.GetByRuleID(ctx, ruleId)
}

func (p *refStore) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
	ruleIds, ok := p.rulesByComponent[componentId]
	if !ok {
		return nil, fmt.Errorf("failed to find rules for component %q", componentId)
	}

	var ruleSets []extensions.RuleSet
	var errs []error
	for ruleId := range ruleIds {
		ruleSet, err := p.GetByRuleID(ctx, ruleId)
		if err != nil {
			errs = append(errs, err)
		}
		ruleSets = append(ruleSets, ruleSet)
	}

	if len(errs) > 0 {
		joinedErr := errors.Join(errs...)
		return ruleSets, fmt.Errorf("failed to find rules for component %q: %w", componentId, joinedErr)
	}

	return ruleSets, nil
}

func (p *refStore) Settings() settings.Settings {
	r := map[string]struct{}{}
	for rule := range p.nodes {
		r[rule] = struct{}{}
	}
	return settings.NewSettings(r, nil)
}

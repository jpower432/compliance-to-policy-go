/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

type PlanRef struct {
	Service  string    `yaml:"service"`
	PluginID plugin.ID `yaml:"pluginID"`
	Plan     *layer4.Layer4
	Loader   Loader
}

type Loader func() (*layer4.Layer4, error)

func (r *PlanRef) Load() error {
	plan, err := r.Loader()
	if err != nil {
		return err
	}
	r.Plan = plan
	return nil
}

// NewContextFromRefs returns an InputContext for the given Layer 3 Policy.
func NewContextFromRefs(refs ...PlanRef) (*InputContext, error) {
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

func newRefStore(refs ...PlanRef) (*refStore, error) {
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

func (p *refStore) indexRef(ref PlanRef) error {
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

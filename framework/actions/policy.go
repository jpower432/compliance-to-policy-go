/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"
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
	return inputCtx, nil
}

var _ rules.Store = (*refStore)(nil)

type refStore struct {
}

func newRefStore(refs ...PlanRef) (*refStore, error) {
	return &refStore{}, nil
}

func (p refStore) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	//TODO implement me
	panic("implement me")
}

func (p refStore) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	//TODO implement me
	panic("implement me")
}

func (p refStore) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
	//TODO implement me
	panic("implement me")
}

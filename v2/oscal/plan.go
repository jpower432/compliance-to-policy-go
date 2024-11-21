package oscal

import (
	"context"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

var _ rules.Store = (*Plan)(nil)

type Plan struct {
	Location string
	store    rules.Store
}

func NewPlan(location string, store rules.Store) *Plan {
	return &Plan{
		Location: location,
		store:    store,
	}
}

func (p *Plan) GetPolicyForComponent(ctx context.Context, componentTitle string) (Policy, error) {
	collectedRules, err := p.store.FindByComponent(ctx, componentTitle)
	if err != nil {
		return Policy{}, err
	}
	// Change for parameter slice
	parameters := make([]extensions.Parameter, 0, len(collectedRules))
	for _, rule := range collectedRules {
		if rule.Rule.Parameter == nil {
			continue
		}
		parameters = append(parameters, *rule.Rule.Parameter)
	}
	policy := Policy{collectedRules, parameters}
	return policy, nil
}

func (p *Plan) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	return p.store.GetByRuleID(ctx, ruleID)
}

func (p *Plan) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	return p.store.GetByCheckID(ctx, checkID)
}

func (p *Plan) FindByComponent(ctx context.Context, componentID string) ([]extensions.RuleSet, error) {
	return p.store.FindByComponent(ctx, componentID)
}

func (p *Plan) All(ctx context.Context) ([]extensions.RuleSet, error) {
	return p.store.All(ctx)
}

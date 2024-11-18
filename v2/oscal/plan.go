package oscal

import (
	"context"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

type Plan struct {
	Location string
	rules.Store
}

func NewPlan(location string, store rules.Store) *Plan {
	return &Plan{
		Location: location,
		Store:    store,
	}
}

func (p *Plan) GetPolicyForComponent(ctx context.Context, componentTitle string) (Policy, error) {
	collectedRules, err := p.FindByComponent(ctx, componentTitle)
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

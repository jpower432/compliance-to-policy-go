package stores

import (
	"context"
	"errors"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var _ rules.Store = (*PlanRefStore)(nil)

type PlanRefStore struct {
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

func NewPlanRefStore(refs ...policy.PlanRef) (*PlanRefStore, error) {
	store := &PlanRefStore{
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

func (p *PlanRefStore) indexRef(ref policy.PlanRef) error {
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
	p.rulesByComponent[ref.PluginID] = ruleIds
	p.rulesByComponent[ref.Service] = ruleIds
	return nil
}

func (p *PlanRefStore) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	ruleSet, ok := p.nodes[ruleID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("rule %q: %w", ruleID, rules.ErrRuleNotFound)
	}
	return ruleSet, nil
}

func (p *PlanRefStore) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	ruleId, ok := p.byCheck[checkID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("failed to find rule for check %q: %w", checkID, rules.ErrRuleNotFound)
	}
	return p.GetByRuleID(ctx, ruleId)
}

func (p *PlanRefStore) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
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

func (p *PlanRefStore) Settings() settings.Settings {
	r := map[string]struct{}{}
	for rule := range p.nodes {
		r[rule] = struct{}{}
	}
	return settings.NewSettings(r, nil)
}

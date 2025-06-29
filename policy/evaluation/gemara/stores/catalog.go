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

var _ rules.Store = (*CatalogRefStore)(nil)

type CatalogRefStore struct {
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

func NewCatalogRefStore(refs ...policy.CatalogRef) (*CatalogRefStore, error) {
	store := &CatalogRefStore{
		nodes:            make(map[string]extensions.RuleSet),
		byCheck:          make(map[string]string),
		rulesByComponent: make(map[string]map[string]struct{}),
	}

	for _, ref := range refs {
		if ref.Catalog == nil {
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

func (p *CatalogRefStore) indexRef(ref policy.CatalogRef) error {
	ruleIds := make(map[string]struct{})
	for _, controlFamilies := range ref.Catalog.ControlFamilies {
		for _, control := range controlFamilies.Controls {
			for _, requirement := range control.AssessmentRequirements {
				ruleSet := extensions.RuleSet{
					Rule: extensions.Rule{
						ID:          requirement.Id,
						Description: requirement.Text,
					},
				}
				p.nodes[requirement.Id] = ruleSet
				ruleIds[requirement.Id] = struct{}{}
			}
		}
	}

	for _, plan := range ref.Plans {
		p.rulesByComponent[plan.Service] = ruleIds
		p.rulesByComponent[plan.PluginID] = ruleIds
	}

	return nil
}

func (p *CatalogRefStore) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	ruleSet, ok := p.nodes[ruleID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("rule %q: %w", ruleID, rules.ErrRuleNotFound)
	}
	return ruleSet, nil
}

func (p *CatalogRefStore) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	ruleId, ok := p.byCheck[checkID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("failed to find rule for check %q: %w", checkID, rules.ErrRuleNotFound)
	}
	return p.GetByRuleID(ctx, ruleId)
}

func (p *CatalogRefStore) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
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

func (p *CatalogRefStore) Settings() settings.Settings {
	r := map[string]struct{}{}
	for rule := range p.nodes {
		r[rule] = struct{}{}
	}
	return settings.NewSettings(r, nil)
}

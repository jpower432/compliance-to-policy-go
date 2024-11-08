package rulefinder

import (
	"fmt"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan"
)

// Naming inspired by https://github.com/oras-project/oras-go/blob/main/internal/graph/memory.go

var _ plan.RuleFinder = &Memory{}

// Memory is a memory based RuleFinder.
type Memory struct {
	ruleSets    map[string]*plan.RuleSet
	byCheck     map[string]*plan.RuleSet
	byComponent map[string][]*plan.RuleSet
	byControl   map[string][]*plan.RuleSet
}

func NewMemoryWithCD(componentDefinition oscalTypes_1_1_2.ComponentDefinition) (*Memory, error) {
	store := Memory{
		ruleSets:    make(map[string]*plan.RuleSet),
		byCheck:     make(map[string]*plan.RuleSet),
		byComponent: make(map[string][]*plan.RuleSet),
		byControl:   make(map[string][]*plan.RuleSet),
	}
	if componentDefinition.Components == nil {
		return nil, fmt.Errorf("no components found")
	}

	componentWideRules := func(component oscalTypes_1_1_2.DefinedComponent) []*plan.RuleSet {
		rules := make([]*plan.RuleSet, 0)
		for _, prop := range *component.Props {
			ruleId := prop.Remarks
			rule, ok := store.ruleSets[ruleId]
			if !ok {
				rule = &plan.RuleSet{}
				store.ruleSets[ruleId] = rule
			}
			switch prop.Name {
			case "Rule_Id":
				rule.RuleID = prop.Value
			case "Rule_Description":
				rule.RuleDescription = prop.Value
			case "Check_Id":
				rule.CheckID = prop.Value
			case "Parameter_Id":
				rule.Parameter.ID = prop.Value
			case "Parameter_Description":
				rule.Parameter.Description = prop.Value
			}
			store.byCheck[rule.CheckID] = rule
			rules = append(rules, rule)
		}
		return rules
	}

	addRuleSet := func(controlId, ruleId string) {
		ruleSet := store.ruleSets[ruleId]
		if _, ok := store.byControl[controlId]; !ok {
			store.byControl[controlId] = []*plan.RuleSet{ruleSet}
		} else {
			store.byControl[controlId] = append(store.byControl[controlId], ruleSet)
		}
	}

	for _, component := range *componentDefinition.Components {
		ruleSets := componentWideRules(component)
		if _, ok := store.byComponent[component.Title]; !ok {
			store.byComponent[component.Title] = ruleSets
		} else {
			store.byComponent[component.Title] = append(store.byComponent[component.Title], ruleSets...)
		}

		if component.ControlImplementations != nil {
			for _, controlImpl := range *component.ControlImplementations {
				for _, implReq := range controlImpl.ImplementedRequirements {

					if implReq.Statements != nil && len(*implReq.Statements) > 0 {
						for _, statement := range *implReq.Statements {
							for _, id := range listRules(statement.Props) {
								addRuleSet(statement.StatementId, id)
							}
						}
					} else {
						for _, id := range listRules(implReq.Props) {
							addRuleSet(implReq.ControlId, id)
						}
					}
				}
			}
		}
	}

	return &store, nil
}

func listRules(props *[]oscalTypes_1_1_2.Property) []string {
	if props == nil {
		return nil
	}
	ruleIds := []string{}
	for _, prop := range *props {
		if prop.Name == "Rule_Id" {
			ruleIds = append(ruleIds, prop.Value)
		}
	}
	return ruleIds
}

func (m Memory) ByCheck(checkId string) (*plan.RuleSet, error) {
	ruleSet, ok := m.byCheck[checkId]
	if !ok {
		return nil, fmt.Errorf("no rule found for check %s", checkId)
	}
	return ruleSet, nil

}

func (m Memory) ByComponent(componentId string) ([]*plan.RuleSet, error) {
	ruleSets, ok := m.byControl[componentId]
	if !ok {
		return nil, fmt.Errorf("no rules found for component %s", componentId)
	}
	return ruleSets, nil
}

func (m Memory) ByControl(controlId string) ([]*plan.RuleSet, error) {
	ruleSets, ok := m.byControl[controlId]
	if !ok {
		return nil, fmt.Errorf("no rules found for control %s", controlId)
	}
	return ruleSets, nil
}

func (m Memory) All() ([]*plan.RuleSet, error) {
	rules := make([]*plan.RuleSet, 0)
	for _, rule := range m.ruleSets {
		rules = append(rules, rule)
	}
	return rules, nil
}

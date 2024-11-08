package plan

import (
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
)

// Finder defines methods for finding information about Rules from different OSCAL Models
type RuleFinder interface {
	ByCheck(string) (*RuleSet, error)
	ByComponent(string) ([]*RuleSet, error)
	ByControl(string) ([]*RuleSet, error)
	All() ([]*RuleSet, error)
}

// RuleSet defines a single rule with all associated metadata
type RuleSet struct {
	// Rule identification
	RuleID string
	// High level description
	RuleDescription string
	// Associated check implementation identification
	CheckID string
	// High level check description
	CheckDescription string
	Parameter        *Parameter
}

// Parameter identifies a parameter or variable that can be used to alter rule logic
type Parameter struct {
	// Parameter Identification
	ID          string
	Description string
	// The selected value for the parameter
	Value string
}

type Policy struct {
	RuleSets   []*RuleSet
	Parameters []*Parameter
}

func (p *Policy) ToProto() *proto.Policy {
	policy := &proto.Policy{
		Rules:      make([]*proto.Rule, 0),
		Parameters: make([]*proto.Parameter, 0),
	}
	for _, rs := range p.RuleSets {
		var parameter *proto.Parameter
		if rs.Parameter != nil {
			parameter = &proto.Parameter{
				Name:          rs.Parameter.ID,
				Description:   rs.Parameter.Description,
				SelectedValue: rs.Parameter.Value,
			}
		}

		ruleSet := &proto.Rule{
			Name:        rs.RuleID,
			Description: rs.RuleDescription,
			Check: &proto.Check{
				Name:        rs.CheckID,
				Description: rs.CheckDescription,
			},
			Parameter: parameter,
		}
		policy.Rules = append(policy.Rules, ruleSet)
		policy.Parameters = append(policy.Parameters, parameter)
	}
	return policy
}

func FromProto(pb *proto.Policy) Policy {
	p := Policy{
		make([]*RuleSet, 0),
		make([]*Parameter, 0),
	}
	for _, r := range pb.Rules {
		var parameter *Parameter
		if r.Parameter != nil {
			parameter = &Parameter{
				ID:          r.Parameter.Name,
				Description: r.Parameter.Description,
				Value:       r.Parameter.SelectedValue,
			}
		}

		ruleSet := &RuleSet{
			RuleID:           r.Name,
			RuleDescription:  r.Description,
			CheckID:          r.Check.Name,
			CheckDescription: r.Check.Description,
			Parameter:        parameter,
		}
		p.RuleSets = append(p.RuleSets, ruleSet)
		p.Parameters = append(p.Parameters, parameter)
	}
	return p
}

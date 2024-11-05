package rules

import (
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
)

type RuleSet struct {
	RuleID           string
	RuleDescription  string
	CheckID          string
	CheckDescription string
	Parameter        *Parameter
}

type Parameter struct {
	ID          string
	Description string
	Value       string
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
		parameter := &proto.Parameter{
			Name:          rs.Parameter.ID,
			Description:   rs.Parameter.Description,
			SelectedValue: rs.Parameter.Value,
		}
		rule := &proto.Rule{
			Name:        rs.RuleID,
			Description: rs.RuleDescription,
			Check: &proto.Check{
				Name:        rs.CheckID,
				Description: rs.CheckDescription,
			},
			Parameter: parameter,
		}
		policy.Rules = append(policy.Rules, rule)
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
		parameter := &Parameter{
			ID:          r.Parameter.Name,
			Description: r.Parameter.Description,
			Value:       r.Parameter.SelectedValue,
		}
		rule := &RuleSet{
			RuleID:           r.Name,
			RuleDescription:  r.Description,
			CheckID:          r.Check.Name,
			CheckDescription: r.Check.Description,
			Parameter:        parameter,
		}
		p.RuleSets = append(p.RuleSets, rule)
		p.Parameters = append(p.Parameters, parameter)
	}
	return p
}

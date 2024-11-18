package oscal

import (
	"fmt"

	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"

	. "github.com/oscal-compass/oscal-sdk-go/extensions"
)

type Policy struct {
	RuleSets   []RuleSet
	Parameters []Parameter
}

func (p *Policy) ToProto() *proto.Policy {
	policy := &proto.Policy{}
	for _, rs := range p.RuleSets {
		var parameter *proto.Parameter
		if rs.Rule.Parameter != nil {
			parameter = &proto.Parameter{
				Name:          rs.Rule.Parameter.ID,
				Description:   rs.Rule.Parameter.Description,
				SelectedValue: rs.Rule.Parameter.Value,
			}
		}

		var checks []*proto.Check
		for _, ch := range rs.Checks {
			fmt.Print(ch.ID)
			check := &proto.Check{
				Name:        ch.ID,
				Description: ch.Description,
			}
			checks = append(checks, check)
		}
		ruleSet := &proto.Rule{
			Name:        rs.Rule.ID,
			Description: rs.Rule.Description,
			Checks:      checks,
			Parameter:   parameter,
		}
		policy.Rules = append(policy.Rules, ruleSet)
		policy.Parameters = append(policy.Parameters, parameter)
	}
	return policy
}

func NewPolicyFromProto(pb *proto.Policy) Policy {
	p := Policy{}
	for _, r := range pb.Rules {
		var parameter Parameter
		if r.Parameter != nil {
			parameter = Parameter{
				ID:          r.Parameter.Name,
				Description: r.Parameter.Description,
				Value:       r.Parameter.SelectedValue,
			}
		}

		var checks []Check
		for _, ch := range r.Checks {
			check := Check{
				ID:          ch.Name,
				Description: ch.Description,
			}
			checks = append(checks, check)
		}

		rule := RuleSet{
			Rule: Rule{
				ID:          r.Name,
				Description: r.Description,
				Parameter:   &parameter,
			},
			Checks: checks,
		}

		p.RuleSets = append(p.RuleSets, rule)
		p.Parameters = append(p.Parameters, parameter)
	}
	return p
}

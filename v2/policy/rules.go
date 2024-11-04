package policy

import (
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
)

type RuleObject struct {
	RuleId               string
	RuleDescription      string
	PolicyId             string
	ParameterId          string
	ParameterDescription string
}

type Policy struct {
	Rules []RuleObject
}

func (p *Policy) ToProto() *proto.Policy {
	return &proto.Policy{}
}

func RulesFromProto(pb *proto.Policy) Policy {
	p := Policy{
		Rules: make([]RuleObject, 0),
	}
	for _, r := range pb.Rules {
		p.Rules = append(p.Rules, RuleObject{
			RuleId:               r.Name,
			RuleDescription:      r.Description,
			PolicyId:             r.Check.Name,
			ParameterId:          r.Parameter.Name,
			ParameterDescription: r.Parameter.Description,
		})
	}
	return p
}

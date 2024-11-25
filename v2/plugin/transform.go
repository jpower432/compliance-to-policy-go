package plugin

import (
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

func PolicyToProto(p *providers.Policy) *proto.Policy {
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

func NewPolicyFromProto(pb *proto.Policy) providers.Policy {
	p := providers.Policy{}
	for _, r := range pb.Rules {
		var parameter extensions.Parameter
		if r.Parameter != nil {
			parameter = extensions.Parameter{
				ID:          r.Parameter.Name,
				Description: r.Parameter.Description,
				Value:       r.Parameter.SelectedValue,
			}
		}

		var checks []extensions.Check
		for _, ch := range r.Checks {
			check := extensions.Check{
				ID:          ch.Name,
				Description: ch.Description,
			}
			checks = append(checks, check)
		}

		rule := extensions.RuleSet{
			Rule: extensions.Rule{
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

var protoByResult = map[providers.Result]proto.Result{
	providers.ResultPass:    proto.Result_RESULT_PASS,
	providers.ResultInvalid: proto.Result_RESULT_UNSPECIFIED,
	providers.ResultError:   proto.Result_RESULT_ERROR,
	providers.ResultWarning: proto.Result_RESULT_WARNING,
	providers.ResultFail:    proto.Result_RESULT_FAILURE,
}

var resultByProto = map[proto.Result]providers.Result{
	proto.Result_RESULT_UNSPECIFIED: providers.ResultInvalid,
	proto.Result_RESULT_ERROR:       providers.ResultError,
	proto.Result_RESULT_WARNING:     providers.ResultWarning,
	proto.Result_RESULT_PASS:        providers.ResultPass,
	proto.Result_RESULT_FAILURE:     providers.ResultFail,
}

func NewResultFromProto(pb *proto.PVPResult) providers.PVPResult {
	result := providers.PVPResult{
		ObservationsByCheck: make([]providers.ObservationByCheck, 0),
	}

	for _, o := range pb.Observations {
		observation := providers.ObservationByCheck{
			Title:       o.Name,
			Description: o.Description,
			Methods:     o.Methods,
			Collected:   o.CollectedAt.AsTime(),
			CheckID:     o.CheckId,
		}
		links := make([]providers.Link, 0)
		for _, ref := range o.EvidenceRefs {
			link := providers.Link{Href: ref}
			links = append(links, link)
		}
		observation.RelevantEvidences = links

		subjects := make([]providers.Subject, 0)
		for _, s := range o.Subjects {
			subject := providers.Subject{
				Title:       s.Title,
				ResourceID:  s.ResourceId,
				Result:      resultByProto[s.Result],
				EvaluatedOn: s.EvaluatedOn.AsTime(),
				Reason:      s.Reason,
			}
			subjects = append(subjects, subject)
		}
		observation.Subjects = subjects
		result.ObservationsByCheck = append(result.ObservationsByCheck, observation)
	}
	return result
}

func ResultsToProto(p *providers.PVPResult) *proto.PVPResult {
	pvpResult := &proto.PVPResult{Observations: make([]*proto.ObservationByCheck, 0)}

	for _, o := range p.ObservationsByCheck {
		observation := &proto.ObservationByCheck{
			Name:        o.Title,
			Description: o.Description,
			CheckId:     o.CheckID,
			Methods:     o.Methods,
			CollectedAt: timestamppb.New(o.Collected),
		}
		subjects := make([]*proto.Subject, 0)
		for _, s := range o.Subjects {
			subject := &proto.Subject{
				Title:       s.Title,
				ResourceId:  s.ResourceID,
				Result:      protoByResult[s.Result],
				EvaluatedOn: timestamppb.New(s.EvaluatedOn),
				Reason:      s.Reason,
			}
			subjects = append(subjects, subject)
		}
		evidences := make([]string, 0)
		for _, evidence := range o.RelevantEvidences {
			evidences = append(evidences, evidence.Href)
		}
		observation.EvidenceRefs = evidences
		observation.Subjects = subjects
		pvpResult.Observations = append(pvpResult.Observations, observation)
	}
	return pvpResult
}

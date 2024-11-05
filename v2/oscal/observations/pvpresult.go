package observations

import proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"

type PVPResult struct {
	ObservationsByCheck []*ObservationByCheck
	Links               []*Link
}

func FromProto(pb *proto.PVPResult) PVPResult {
	result := PVPResult{
		ObservationsByCheck: make([]*ObservationByCheck, 0),
	}

	for _, o := range pb.Observations {
		observation := ObservationByCheck{
			Title:       o.Name,
			Description: o.Description,
			Methods:     o.Methods,
		}
		links := make([]*Link, 0)
		for _, ref := range o.EvidenceRefs {
			link := &Link{Href: ref}
			links = append(links, link)
		}
		observation.RelevantEvidences = links
		subjects := make([]*Subject, 0)
		for _, s := range o.Subjects {
			subject := &Subject{
				Title:       s.Title,
				ResourceID:  s.ResourceId,
				Result:      resultByProto[s.Result],
				EvaluatedOn: s.EvaluatedOn.AsTime(),
				Reason:      s.Reason,
			}
			subjects = append(subjects, subject)
		}
		observation.Subjects = subjects
		result.ObservationsByCheck = append(result.ObservationsByCheck, &observation)
	}
	return result
}

func (p *PVPResult) ToProto() *proto.PVPResult {
	pvpResult := &proto.PVPResult{Observations: make([]*proto.ObservationByCheck, 0)}

	for _, o := range p.ObservationsByCheck {
		observation := &proto.ObservationByCheck{
			Name:        o.Title,
			Description: o.Description,
			CheckId:     o.CheckID,
		}
		subjects := make([]*proto.Subject, 0)
		for _, s := range o.Subjects {
			subject := &proto.Subject{
				Title:      s.Title,
				ResourceId: s.ResourceID,
				Result:     protoByResult[s.Result],
				Reason:     s.Reason,
			}
			subjects = append(subjects, subject)
		}
		observation.Subjects = subjects
		pvpResult.Observations = append(pvpResult.Observations, observation)
	}
	return pvpResult
}

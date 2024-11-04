package policy

import (
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"time"
)

type Prop struct {
	Name    string
	Ns      string
	Value   string
	Class   string
	Remarks string
}

type Link struct {
	Href             string
	Rel              string
	MediaType        string
	ResourceFragment string
	Text             string
}

type Subject struct {
	Type  string
	Title string
	Props []Prop
}

type ObservationByCheck struct {
	Title            string
	Description      string
	CheckId          string
	Methods          []string
	Types            []string
	Subjects         []Subject
	RelevantEvidence []string
	Collected        time.Time
	Expires          time.Time
}

type PVPResult struct {
	Observations []ObservationByCheck
}

func (p *PVPResult) ToProto() *proto.PVPResult {
	pvpResult := &proto.PVPResult{Observations: make([]*proto.ObservationByCheck, 0)}

	for _, observation := range p.Observations {
		o := &proto.ObservationByCheck{
			Name:         observation.Title,
			Description:  observation.Description,
			CheckId:      observation.CheckId,
			EvidenceRefs: observation.RelevantEvidence,
		}
		pvpResult.Observations = append(pvpResult.Observations, o)
	}
	return pvpResult
}

func ResultFromProto(pb *proto.PVPResult) PVPResult {
	result := PVPResult{
		Observations: make([]ObservationByCheck, len(pb.Observations)),
	}
	for i, o := range pb.Observations {
		result.Observations[i] = ObservationByCheck{
			Title:            o.Name,
			Description:      o.Description,
			Methods:          o.Methods,
			RelevantEvidence: o.EvidenceRefs,
		}
	}
	return result
}

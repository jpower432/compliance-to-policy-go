package report

import (
	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/observations"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/rules"
)

type Reporter struct {
	plan plan.Decomposer
}

func New(plan plan.Decomposer) *Reporter {
	return &Reporter{
		plan: plan,
	}
}

func (r *Reporter) ToOSCAL(results observations.PVPResult) (oscalTypes.AssessmentResults, error) {
	arResult := oscalTypes.AssessmentResults{
		ImportAp: oscalTypes.ImportAp{
			Href: r.plan.Location(),
		},
	}
	oscalObservations := make([]oscalTypes.Observation, 0)

	for _, observation := range results.ObservationsByCheck {
		ruleSet, err := r.plan.RuleSetByCheck(observation.CheckID)
		if err != nil {
			return arResult, err
		}
		oscalObservations = append(oscalObservations, r.toOSCALObservation(*observation, ruleSet))
	}

	arResult.Results = []oscalTypes.Result{
		{
			Observations: &oscalObservations,
		},
	}
	return arResult, nil
}

func (r *Reporter) toOSCALObservation(observationByCheck observations.ObservationByCheck, ruleSet rules.RuleSet) oscalTypes.Observation {
	subjects := make([]oscalTypes.SubjectReference, 0)
	for _, subject := range observationByCheck.Subjects {
		props := []oscalTypes.Property{
			{
				Name:  "resource-id",
				Value: subject.ResourceID,
			},
			{
				Name:  "result",
				Value: subject.Result.String(),
			},
			{
				Name:  "evaluated-on",
				Value: subject.EvaluatedOn.String(),
			},
			{
				Name:  "reason",
				Value: subject.Reason,
			},
		}

		s := oscalTypes.SubjectReference{
			SubjectUuid: uuid.NewUUID(),
			Title:       subject.Title,
			Type:        subject.Type,
			Props:       &props,
		}
		subjects = append(subjects, s)
	}

	relevantEvidences := make([]oscalTypes.RelevantEvidence, 0)
	if observationByCheck.RelevantEvidences != nil {
		for _, rel := range observationByCheck.RelevantEvidences {
			re := oscalTypes.RelevantEvidence{
				Href:        rel.Href,
				Description: rel.Description,
			}
			relevantEvidences = append(relevantEvidences, re)
		}
	}

	props := []oscalTypes.Property{
		{
			Name:  "assessment-rule-id",
			Value: ruleSet.RuleID,
		},
	}

	observation := oscalTypes.Observation{
		UUID:        uuid.NewUUID(),
		Title:       observationByCheck.Title,
		Description: observationByCheck.Description,
		Methods:     observationByCheck.Methods,
		Props:       &props,
		Subjects:    &subjects,
		Collected:   observationByCheck.Collected,
	}
	if len(relevantEvidences) > 0 {
		observation.RelevantEvidence = &relevantEvidences
	}

	return observation
}

func (r *Reporter) ToMarkdown(results oscalTypes.AssessmentResults) error {
	return nil
}

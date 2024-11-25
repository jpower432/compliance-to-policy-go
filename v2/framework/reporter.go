package framework

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/oscal-sdk-go/rules"

	"github.com/oscal-compass/oscal-sdk-go/extensions"

	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

type Reporter struct {
	store rules.Store
}

func NewReporter(store rules.Store) *Reporter {
	return &Reporter{
		store: store,
	}
}

func (r *Reporter) ToOSCAL(ctx context.Context, planHref string, results []providers.PVPResult) (oscalTypes.AssessmentResults, error) {
	arResult := oscalTypes.AssessmentResults{
		ImportAp: oscalTypes.ImportAp{
			Href: planHref,
		},
	}
	oscalObservations := make([]oscalTypes.Observation, 0)

	for _, result := range results {
		for _, observation := range result.ObservationsByCheck {
			rule, err := r.store.GetByCheckID(ctx, observation.CheckID)
			if err != nil {
				if !errors.Is(err, rules.ErrRuleNotFound) {
					return arResult, fmt.Errorf("failed to convert observation for check: %w", err)
				}
				log.Print(err)
			}
			oscalObservations = append(oscalObservations, r.toOSCALObservation(observation, rule))
		}
	}

	arResult.Results = []oscalTypes.Result{
		{
			Observations: &oscalObservations,
		},
	}
	return arResult, nil
}

func (r *Reporter) toOSCALObservation(observationByCheck providers.ObservationByCheck, ruleSet extensions.RuleSet) oscalTypes.Observation {
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
			Value: ruleSet.Rule.ID,
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

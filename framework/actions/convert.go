package actions

import (
	"context"
	"time"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// For reporting in OSCAL ecosystem

func ObservationsFromEvaluation(eval layer4.Layer4) policy.PVPResult {
	result := policy.PVPResult{}
	for _, controlEval := range eval.ControlEvaluations {
		selection := oscalTypes.AssessedControlsSelectControlById{
			ControlId:    controlEval.ControlID,
			StatementIds: &[]string{},
		}
		for _, assessment := range controlEval.Assessments {
			*selection.StatementIds = append(*selection.StatementIds, assessment.RequirementID)
			for _, method := range assessment.Methods {
				obs := observation(method, eval.EndTime, assessment.RequirementID)
				result.ObservationsByCheck = append(result.ObservationsByCheck, obs)
			}
		}
	}
	return result
}

func observation(method layer4.AssessmentMethod, end time.Time, req string) policy.ObservationByCheck {
	result := policy.ResultInvalid
	if method.Result != nil && method.Result.Status != "" {
		switch method.Result.Status {
		case "passed":
			result = policy.ResultPass
		case "failed":
			result = policy.ResultFail
		case "error":
			result = policy.ResultError
		}
	}
	if method.RemediationGuide == "" {
		method.RemediationGuide = "N/A"
	}

	return policy.ObservationByCheck{
		Collected:   end,
		Description: method.Description,
		Title:       method.Name,
		CheckID:     method.Name,
		Requirement: req,
		Methods: []string{
			"TEST",
		},
		// FIXME: This would need to be filled out by more granular evidence information
		// in a L4 evaluation.
		Subjects: []policy.Subject{
			{
				Title:       "",
				Type:        "resource",
				ResourceID:  "",
				Result:      result,
				EvaluatedOn: time.Now(),
				Reason:      "",
			},
		},
	}
}

// To produce evaluations for Gemara ecosystem

func Layer4FromResults(ctx context.Context, inputContext *InputContext, catalogID string, results []policy.PVPResult) ([]policy.PlanRef, error) {
	var planRef []policy.PlanRef
	for _, result := range results {
		ref := policy.PlanRef{
			Plan: &layer4.Layer4{},
		}
		ref.Plan.CatalogID = catalogID
		ref.Plan.StartTime = time.Now()
		checksByRule := make(map[string][]policy.ObservationByCheck)
		store := inputContext.Store()
		for _, observationByCheck := range result.ObservationsByCheck {
			rule, err := store.GetByRuleID(ctx, observationByCheck.Requirement)
			if err != nil {
				return planRef, err
			}
			checksByRule[rule.Rule.ID] = append(checksByRule[rule.Rule.ID], observationByCheck)
		}

		for _, controlEvals := range ref.Plan.ControlEvaluations {
			for i := range controlEvals.Assessments {
				checks := checksByRule[controlEvals.Assessments[i].RequirementID]
				controlEvals.Assessments[i].Methods = getMethods(checks)
			}
		}
		ref.Plan.EndTime = time.Now()
		planRef = append(planRef, ref)
	}
	return planRef, nil
}

func getMethods(assessmentMethods []policy.ObservationByCheck) []layer4.AssessmentMethod {
	var methods []layer4.AssessmentMethod
	for _, method := range assessmentMethods {
		l4Assessment := layer4.AssessmentMethod{
			Name:        method.Title,
			Description: method.Description,
			Run:         true,
			Result: &layer4.AssessmentResult{
				Status: layer4.Status(method.Subjects[0].Result.String()),
			},
		}
		methods = append(methods, l4Assessment)
	}

	return methods
}

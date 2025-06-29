package gemara

import (
	"context"
	"fmt"
	"time"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
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

func Layer4FromResults(ctx context.Context, inputContext *evaluation.InputContext, results []policy.PVPResult, catalog *policy.CatalogRef) ([]policy.PlanRef, error) {
	checksByRule := make(map[string][]policy.ObservationByCheck)
	store := inputContext.Store()
	for _, result := range results {
		for _, observationByCheck := range result.ObservationsByCheck {
			rule, err := store.GetByCheckID(ctx, observationByCheck.CheckID)
			if err != nil {
				return nil, err
			}
			checksByRule[rule.Rule.ID] = append(checksByRule[rule.Rule.ID], observationByCheck)
		}
	}

	var plans []policy.PlanRef
	for _, planRef := range catalog.Plans {
		planRef.Plan = layer4.NewEvaluation(*catalog.Catalog)
		planRef.Plan.StartTime = time.Now()
		for i := range planRef.Plan.ControlEvaluations {
			for j := range planRef.Plan.ControlEvaluations[i].Assessments {
				checks := checksByRule[planRef.Plan.ControlEvaluations[i].Assessments[j].RequirementID]
				planRef.Plan.ControlEvaluations[i].Assessments[j].Methods = getMethods(checks)
			}
		}
		planRef.Plan.EndTime = time.Now()
		plans = append(plans, planRef)
	}

	return plans, nil
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

type Marshalable struct {
	layer4.Layer4
}

// Temporary until unmarshalling is done upstream

func (mc Marshalable) MarshalYAML() (interface{}, error) {
	outputMap := make(map[string]interface{})
	outputMap["catalog_id"] = mc.CatalogID
	outputMap["start_time"] = mc.StartTime
	outputMap["end_time"] = mc.EndTime
	outputMap["corrupted_state"] = mc.CorruptedState

	controlEvals := []map[string]interface{}{}
	for _, controlEval := range mc.ControlEvaluations {
		evalMap := make(map[string]interface{})
		evalMap["control_id"] = controlEval.ControlID
		assessments := []map[string]interface{}{}
		for _, assessment := range controlEval.Assessments {
			assessmentMap := make(map[string]interface{})
			assessmentMap["requirement_id"] = assessment.RequirementID
			methods := []map[string]interface{}{}
			for _, method := range assessment.Methods {
				methodMap := make(map[string]interface{})
				methodMap["name"] = method.Name
				methodMap["description"] = method.Description
				methodMap["run"] = method.Run
				if method.Result != nil {
					methodMap["result"] = map[string]interface{}{
						"status": method.Result.Status,
					}
				}
				methods = append(methods, methodMap)
			}
			assessmentMap["methods"] = methods
			assessments = append(assessments, assessmentMap)
		}
		evalMap["assessments"] = assessments
		controlEvals = append(controlEvals, evalMap)
	}
	outputMap["evaluations"] = controlEvals

	// Return the map, which yaml.Marshal will then convert into YAML.
	return outputMap, nil
}

// MarshalConfigWithFunction takes a Config struct and returns its YAML representation
// by leveraging the custom MarshalYAML implementation.
func MarshalConfigWithFunction(eval *layer4.Layer4) ([]byte, error) {
	marshalable := Marshalable{*eval}

	yamlData, err := yaml.Marshal(marshalable)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config to YAML using custom marshaler: %w", err)
	}
	return yamlData, nil
}

/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"

	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/revanite-io/sci/layer2"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

func Evaluate(ctx context.Context, inputContext *InputContext, controls []layer2.Control, results []policy.PVPResult) ([]layer4.Evaluation, error) {
	var evals []layer4.Evaluation
	for _, result := range results {
		eval, err := createEvaluation(ctx, inputContext.Store(), controls, result)
		if err != nil {
			return nil, err
		}
		evals = append(evals, eval)
	}
	return evals, nil
}

func createEvaluation(ctx context.Context, store rules.Store, controls []layer2.Control, result policy.PVPResult) (layer4.Evaluation, error) {
	checksByRule := make(map[string][]policy.ObservationByCheck)
	for _, observationByCheck := range result.ObservationsByCheck {
		rule, err := store.GetByCheckID(ctx, observationByCheck.CheckID)
		if err != nil {
			return layer4.Evaluation{}, err
		}
		checksByRule[rule.Rule.ID] = append(checksByRule[rule.Rule.ID], observationByCheck)
	}

	// Assuming rule will align with the requirement id
	eval := layer4.Evaluation{}
	for _, control := range controls {
		controlEval := layer4.ControlEvaluation{
			ControlID: control.Id,
		}
		for _, req := range control.AssessmentRequirements {
			assessment := layer4.ForControlRequirement(control.Id, req)
			checks := checksByRule[req.Id]
			assessment.Methods = getMethods(checks)
		}
		eval.ControlEvaluations = append(eval.ControlEvaluations, controlEval)
	}
	return eval, nil
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
			// Look into "MethodExecutor"

		}
		methods = append(methods, l4Assessment)
	}

	return methods
}

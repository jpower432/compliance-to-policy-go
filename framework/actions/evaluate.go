/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/revanite-io/sci/layer4"

	fpolicy "github.com/oscal-compass/compliance-to-policy-go/v2/framework/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// Evaluate updates the given Layer 4 evaluation based on PVP Results.
func Evaluate(ctx context.Context, inputContext *InputContext, ref *fpolicy.PlanRef, provider policy.Provider) error {
	if ref.Plan == nil {
		err := ref.Load()
		if err != nil {
			return err
		}
	}

	ref.Plan.StartTime = time.Now()

	appliedRuleSet, err := settings.ApplyToComponent(ctx, ref.PluginID.String(), inputContext.Store(), inputContext.Settings)
	if err != nil {
		return fmt.Errorf("failed to get rule sets for component %s: %w", ref.Service, err)
	}
	results, err := provider.GetResults(appliedRuleSet)
	if err != nil {
		return fmt.Errorf("plugin %v: %w", ref.PluginID, err)
	}

	checksByRule := make(map[string][]policy.ObservationByCheck)
	store := inputContext.Store()
	for _, observationByCheck := range results.ObservationsByCheck {
		rule, err := store.GetByCheckID(ctx, observationByCheck.CheckID)
		if err != nil {
			return err
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
	return nil
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

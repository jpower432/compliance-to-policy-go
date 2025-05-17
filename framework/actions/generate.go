/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/revanite-io/sci/layer2"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/logging"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// GeneratePolicy action identifies policy configuration for each provider in the given pluginSet to execute the Generate() method
// each policy.Provider.
//
// The rule set passed to each plugin can be configured with compliance specific settings based on the InputContext.
func GeneratePolicy(ctx context.Context, inputContext *InputContext, pluginSet map[plugin.ID]policy.Provider) error {
	log := logging.GetLogger("generator")

	for providerId, policyPlugin := range pluginSet {
		componentTitle, err := inputContext.ProviderTitle(providerId)
		if err != nil {
			if errors.Is(err, ErrMissingProvider) {
				log.Warn(fmt.Sprintf("skipping %s provider: missing validation component", providerId))
				continue
			}
			return err
		}
		log.Debug(fmt.Sprintf("Generating policy for provider %s", providerId))

		appliedRuleSet, err := settings.ApplyToComponent(ctx, componentTitle, inputContext.Store(), inputContext.Settings)
		if err != nil {
			return fmt.Errorf("failed to get rule sets for component %s: %w", componentTitle, err)
		}
		if err := policyPlugin.Generate(appliedRuleSet); err != nil {
			return fmt.Errorf("plugin %s: %w", providerId, err)
		}
	}
	return nil
}

// GenerateEvaluation returns a Layer4 evaluation plan based on a Layer 2 catalog and action context.
// This should also generate policy.
func GenerateEvaluation(catalog layer2.Layer2, provider policy.Provider) (*layer4.Layer4, error) {
	var appliedRuleSets policy.Policy

	evaluation := layer4.NewEvaluation(catalog)
	for _, controlEval := range evaluation.ControlEvaluations {
		for _, assessment := range controlEval.Assessments {
			ruleSet := extensions.RuleSet{
				Rule: extensions.Rule{
					ID: assessment.RequirementID,
				},
			}
			for _, method := range assessment.Methods {
				check := extensions.Check{
					ID:          method.Name,
					Description: method.Description,
				}
				ruleSet.Checks = append(ruleSet.Checks, check)
			}
			appliedRuleSets = append(appliedRuleSets, ruleSet)
		}
	}

	if err := provider.Generate(appliedRuleSets); err != nil {
		return nil, err
	}
	return evaluation, nil
}

package framework

import (
	"context"
	"errors"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/config"
	"github.com/oscal-compass/compliance-to-policy-go/v2/logging"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// GeneratePolicy identifies policy configuration for each provider in the given pluginSet to execute the Generate() method
// each policy.Provider. The rule set passed to each plugin can be configured with compliance specific settings with the
// complianceSettings input.
func GeneratePolicy(ctx context.Context, pluginSet map[string]policy.Generator, target *config.Target) error {
	log := logging.GetLogger("aggregator")

	for providerId, policyPlugin := range pluginSet {
		componentTitle, err := target.PluginTitle(providerId)
		if err != nil {
			if errors.Is(err, config.ErrMissingProvider) {
				log.Warn(fmt.Sprintf("skipping %s provider: missing validation component", providerId))
				continue
			}
			return err
		}
		log.Debug(fmt.Sprintf("Generating policy for provider %s", providerId))

		appliedRuleSet, err := settings.ApplyToComponent(ctx, componentTitle, target.Store(), target.Settings)
		if err != nil {
			return fmt.Errorf("failed to get rule sets for component %s: %w", componentTitle, err)
		}
		if err := policyPlugin.Generate(appliedRuleSet); err != nil {
			return fmt.Errorf("plugin %s: %w", providerId, err)
		}
	}
	return nil
}

// AggregateResults identifies policy configuration for each provider in the given pluginSet to execute the GetResults() method
// each policy.Provider. The rule set passed to each plugin can be configured with compliance specific settings with the
// complianceSettings input.
func AggregateResults(ctx context.Context, pluginSet map[string]policy.Aggregator, target *config.Target) ([]policy.PVPResult, error) {
	var allResults []policy.PVPResult
	log := logging.GetLogger("aggregator")
	for providerId, policyPlugin := range pluginSet {
		componentTitle, err := target.PluginTitle(providerId)
		if err != nil {
			return nil, err
		}
		log.Debug(fmt.Sprintf("Aggregating results for provider %s", providerId))
		appliedRuleSet, err := settings.ApplyToComponent(ctx, componentTitle, target.Store(), target.Settings)
		if err != nil {
			return allResults, fmt.Errorf("failed to get rule sets for component %s: %w", componentTitle, err)
		}

		pluginResults, err := policyPlugin.GetResults(appliedRuleSet)
		if err != nil {
			return allResults, fmt.Errorf("plugin %s: %w", providerId, err)
		}
		allResults = append(allResults, pluginResults)
	}
	return allResults, nil
}

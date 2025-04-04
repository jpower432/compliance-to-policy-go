package action

import (
	"context"
	"errors"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/logging"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// GeneratePolicy identifies policy configuration for each provider in the given pluginSet to execute the Generate() method
// each policy.Provider. The rule set passed to each plugin can be configured with compliance specific settings with the
// complianceSettings input.
func GeneratePolicy(ctx context.Context, pluginSet map[string]policy.Generator, target *Target) error {
	log := logging.GetLogger("aggregator")

	for providerId, policyPlugin := range pluginSet {
		componentTitle, err := target.PluginTitle(providerId)
		if err != nil {
			if errors.Is(err, ErrMissingProvider) {
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

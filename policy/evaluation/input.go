/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package evaluation

import (
	"errors"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

var ErrMissingProvider = errors.New("missing title for provider")

// InputContext is used to configure action behavior from parsed OSCAL documents.
type InputContext struct {
	// requestedProviders stores resolved provider IDs to the original component title from
	// parsed OSCAL Components.
	//
	// The original component title is needed to get information for the validation
	// component in the rules.Store (which provides input for the corresponding policy.Provider
	// id).
	requestedProviders map[plugin.ID]string
	// rulesStore contains indexed information about parsed RuleSets
	// which can be searched for the corresponding component title.
	rulesStore rules.Store
	// Settings define adjustable rule settings parsed from framework-specific implementation
	Settings settings.Settings
	// Applicability defines the profiles for a group a checks
	Applicability string
}

func NewContext(providers map[plugin.ID]string, store rules.Store) *InputContext {
	return &InputContext{
		requestedProviders: providers,
		rulesStore:         store,
	}
}

// RequestedProviders returns the provider ids requested in the parsed input.
func (t *InputContext) RequestedProviders() []plugin.ID {
	var requestedIds []plugin.ID
	for id := range t.requestedProviders {
		requestedIds = append(requestedIds, id)
	}
	return requestedIds
}

// ProviderTitle return the original component Title for the given provider id.
func (t *InputContext) ProviderTitle(providerId plugin.ID) (string, error) {
	title, ok := t.requestedProviders[providerId]
	if !ok {
		return "", fmt.Errorf("%s:%w", providerId, ErrMissingProvider)
	}
	return title, nil
}

// Store returns the underlying rules.Store with indexed RuleSets.
func (t *InputContext) Store() rules.Store {
	return t.rulesStore
}

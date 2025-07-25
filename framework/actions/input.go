/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

const pluginComponentType = "validation"

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
	// action concurrency
	MaxConcurrency int
}

// NewContext returns an InputContext for the given OSCAL Components.
func NewContext(components []components.Component) (*InputContext, error) {
	inputCtx := &InputContext{
		requestedProviders: make(map[plugin.ID]string),
		MaxConcurrency:     3,
	}
	for _, comp := range components {
		if comp.Type() == pluginComponentType {
			pluginId, err := GetPluginIDFromComponent(comp)
			if err != nil {
				return inputCtx, err
			}
			inputCtx.requestedProviders[pluginId] = comp.Title()
		}
	}
	store, err := DefaultStore(components)
	if err != nil {
		return inputCtx, err
	}
	inputCtx.rulesStore = store
	return inputCtx, nil
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

// DefaultStore returns a default rules.MemoryStore with indexed information from the given
// Components.
func DefaultStore(allComponents []components.Component) (*rules.MemoryStore, error) {
	store := rules.NewMemoryStore()
	err := store.IndexAll(allComponents)
	if err != nil {
		return store, err
	}
	return store, nil
}

// GetPluginIDFromComponent returns the normalized plugin identifier defined by the OSCAL Component
// of type "validation".
func GetPluginIDFromComponent(component components.Component) (plugin.ID, error) {
	title := strings.TrimSpace(component.Title())
	if title == "" {
		return "", fmt.Errorf("component is missing a title")
	}

	title = strings.ToLower(title)
	id := plugin.ID(title)
	if !id.Validate() {
		return "", fmt.Errorf("invalid plugin id %s", title)
	}
	return id, nil
}

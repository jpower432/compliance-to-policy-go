/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package action

import (
	"errors"
	"fmt"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

const pluginComponentType = "validation"

var ErrMissingProvider = errors.New("missing title for provider")

type Target struct {
	// rulesStore contains indexed information about available RuleSets
	// which can be searched for the component title.
	rulesStore rules.Store
	// pluginIdMap stores resolved plugin IDs to the original component title from the
	// OSCAL Component Definitions.
	//
	// The original component title is needed to get information for the validation
	// component in the rules.Store (which provides input for the corresponding policy.Provider
	// plugin).
	requestedPlugins map[string]string
	// Adjustable target settings
	Settings settings.Settings
}

func NewTargetFromComponentDefinition(compDef oscalTypes.ComponentDefinition) (*Target, error) {
	target := &Target{
		requestedPlugins: make(map[string]string),
	}
	var allComponents []components.Component
	if compDef.Components == nil {
		return target, fmt.Errorf("components cannot be empty")
	}
	for _, component := range *compDef.Components {
		compAdapter := components.NewDefinedComponentAdapter(component)
		if compAdapter.Type() == pluginComponentType {
			pluginId, err := GetPluginIDFromComponent(compAdapter)
			if err != nil {
				return target, err
			}
			target.requestedPlugins[pluginId] = component.Title
		}
		allComponents = append(allComponents, compAdapter)
	}
	store, err := DefaultStore(allComponents)
	if err != nil {
		return target, err
	}
	target.rulesStore = store
	return target, nil
}

func (t *Target) RequiredPlugins() []string {
	var requestedIds []string
	for id := range t.requestedPlugins {
		requestedIds = append(requestedIds, id)
	}
	return requestedIds
}

func (t *Target) PluginTitle(providerId string) (string, error) {
	title, ok := t.requestedPlugins[providerId]
	if !ok {
		return "", fmt.Errorf("%s:%w", providerId, ErrMissingProvider)
	}
	return title, nil
}

func (t *Target) Store() rules.Store {
	return t.rulesStore
}

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
func GetPluginIDFromComponent(component components.Component) (string, error) {
	title := strings.TrimSpace(component.Title())
	if title == "" {
		return "", fmt.Errorf("component is missing a title")
	}

	title = strings.ToLower(title)
	if !plugin.IdentifierPattern.MatchString(title) {
		return "", fmt.Errorf("invalid plugin id %s", title)
	}
	return title, nil
}

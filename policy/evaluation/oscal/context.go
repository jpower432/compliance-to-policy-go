package oscal

import (
	"errors"
	"fmt"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/settings"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
)

const pluginComponentType = "validation"

func NewContext(ap *oscalTypes.AssessmentPlan) (*evaluation.InputContext, error) {
	if ap.LocalDefinitions == nil || ap.LocalDefinitions.Activities == nil || ap.AssessmentAssets.Components == nil {
		return nil, errors.New("error converting component definition to assessment plan")
	}

	var allComponents []components.Component
	for _, component := range *ap.AssessmentAssets.Components {
		compAdapter := components.NewSystemComponentAdapter(component)
		allComponents = append(allComponents, compAdapter)
	}

	inputCtx, err := newContext(allComponents)
	if err != nil {
		return nil, err
	}

	apSettings := settings.NewAssessmentActivitiesSettings(*ap.LocalDefinitions.Activities)
	inputCtx.Settings = apSettings

	return inputCtx, nil
}

// NewContext returns an InputContext for the given OSCAL Components.
func newContext(components []components.Component) (*evaluation.InputContext, error) {
	requestedProviders := make(map[plugin.ID]string)
	for _, comp := range components {
		if comp.Type() == pluginComponentType {
			pluginId, err := GetPluginIDFromComponent(comp)
			if err != nil {
				return nil, err
			}
			requestedProviders[pluginId] = comp.Title()
		}
	}
	store, err := DefaultStore(components)
	if err != nil {
		return nil, err
	}
	inputCtx := evaluation.NewContext(requestedProviders, store)
	return inputCtx, nil
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

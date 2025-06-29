package gemara

import (
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation/gemara/stores"
)

// NewContextFromPlanRefs returns an InputContext for evaluations.
func NewContextFromPlanRefs(refs ...policy.PlanRef) (*evaluation.InputContext, error) {
	requestedProviders := make(map[plugin.ID]string)
	for _, ref := range refs {
		requestedProviders[plugin.ID(ref.PluginID)] = ref.PluginID
	}
	store, err := stores.NewPlanRefStore(refs...)
	if err != nil {
		return nil, err
	}

	inputCtx := evaluation.NewContext(requestedProviders, store)

	inputCtx.Settings = store.Settings()
	return inputCtx, nil
}

// NewContextFromCatalogRefs returns an InputContext for catalogs.
func NewContextFromCatalogRefs(refs ...policy.CatalogRef) (*evaluation.InputContext, error) {
	requestedProviders := make(map[plugin.ID]string)

	for _, ref := range refs {
		for _, plan := range ref.Plans {
			requestedProviders[plugin.ID(plan.PluginID)] = plan.PluginID
		}
	}
	store, err := stores.NewCatalogRefStore(refs...)
	if err != nil {
		return nil, err
	}
	inputCtx := evaluation.NewContext(requestedProviders, store)
	inputCtx.Settings = store.Settings()
	return inputCtx, nil
}

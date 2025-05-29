package policy

import (
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

// Policy is a WIP policy implementation while SCI Layer 3 is under development.
// For prototyping there is 1:1 relationship
// between the target and validator.
type Policy struct {
	Catalogs []string  `yaml:"catalogs"`
	Refs     []PlanRef `yaml:"services"`
}

type PlanRef struct {
	Service  string    `yaml:"service"`
	PluginID plugin.ID `yaml:"pluginID"`
	Plan     *layer4.Layer4
	// Lazy loading
	Loader Loader
}

type Loader func() (*layer4.Layer4, error)

func (r *PlanRef) Load() error {
	plan, err := r.Loader()
	if err != nil {
		return err
	}
	r.Plan = plan
	return nil
}

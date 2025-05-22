package agentkit

import (
	"context"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/resource"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

type Agent struct {
	provider policy.Provider
	plan     actions.PlanRef
}

// NewAgent creates a new C2P Agent.
func NewAgent(provider policy.Provider, plan actions.PlanRef) *Agent {
	agent := &Agent{provider: provider, plan: plan}
	return agent
}

type runOptions struct {
	exportURL string
}

func (o *runOptions) defaults() {
	o.exportURL = "localhost:8080"
}

type RunOption func(ro *runOptions)

// Perhaps set the exporter object instead?

func RunWithExporterURL(url string) RunOption {
	return func(ro *runOptions) {
		ro.exportURL = url
	}
}

func (a *Agent) Run(ctx context.Context, opts ...RunOption) error {
	options := runOptions{}
	options.defaults()
	for _, opt := range opts {
		opt(&options)
	}

	inputContext, err := actions.NewContextFromRefs(a.plan)
	if err != nil {
		return err
	}

	rs, err := actions.Evaluate(ctx, inputContext, a.plan, a.provider)
	if err != nil {
		return err
	}

	artifact := resource.NewAttestation(options.exportURL)
	err = artifact.Attach(rs, *a.plan.Plan)
	if err != nil {
		return err
	}

	err = artifact.Export(ctx)
	if err != nil {
		return err
	}
	return nil
}

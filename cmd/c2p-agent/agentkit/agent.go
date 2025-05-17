package agentkit

import (
	"context"
	"errors"
	"sync"

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

	var wg sync.WaitGroup
	errs := make(chan error)
	errsDone := make(chan struct{})

	var resultErrs []error
	go func() {
		for err := range errs {
			resultErrs = append(resultErrs, err)
		}
		errsDone <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		inputContext, err := actions.NewContextFromRefs(a.plan)
		if err != nil {
			errs <- err
		}

		rs, err := actions.Evaluate(ctx, inputContext, a.plan, a.provider)
		if err != nil {
			errs <- err
		}

		artifact := resource.NewAttestation(options.exportURL)
		err = artifact.Attach(rs, *a.plan.Plan)
		if err != nil {
			errs <- err
		}

		err = artifact.Export(ctx)
		if err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	<-errsDone

	return errors.Join(resultErrs...)
}

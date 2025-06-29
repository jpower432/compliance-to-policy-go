package pipelines

import (
	"context"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/actions"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
)

// Pipeline defines components involves in an evaluation pipeline
type Pipeline struct {
	// Planner provides evaluation plans
	Planner    evaluation.Provider
	ReportPath string
}

func (p *Pipeline) GenerateWorkflow(ctx context.Context, frameworkConfig *framework.C2PConfig, configSelections framework.PluginConfig) error {
	inputContext, err := p.Planner.Plan()
	if err != nil {
		return err
	}

	manager, err := framework.NewPluginManager(frameworkConfig)
	if err != nil {
		return err
	}
	foundPlugins, err := manager.FindRequestedPlugins(inputContext.RequestedProviders())
	if err != nil {
		return err
	}

	launchedPlugins, err := manager.LaunchPolicyPlugins(foundPlugins, configSelections)
	// Defer clean before returning an error to avoid unterminated processes
	defer manager.Clean()
	if err != nil {
		return err
	}

	return actions.GeneratePolicy(ctx, inputContext, launchedPlugins)
}

func (p *Pipeline) ResultsPipeline(ctx context.Context, frameworkConfig *framework.C2PConfig, configSelections framework.PluginConfig) error {
	inputContext, err := p.Planner.Plan()
	if err != nil {
		return err
	}

	manager, err := framework.NewPluginManager(frameworkConfig)
	if err != nil {
		return err
	}
	foundPlugins, err := manager.FindRequestedPlugins(inputContext.RequestedProviders())
	if err != nil {
		return err
	}

	launchedPlugins, err := manager.LaunchPolicyPlugins(foundPlugins, configSelections)
	// Defer clean before returning an error to avoid unterminated processes
	defer manager.Clean()
	if err != nil {
		return err
	}

	results, err := actions.AggregateResults(ctx, inputContext, launchedPlugins)
	if err != nil {
		return err
	}

	return p.Planner.Report(ctx, inputContext, p.ReportPath, results)
}

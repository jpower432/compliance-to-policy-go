package subcommands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-hclog"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/revanite-io/sci/layer4"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation/gemara"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation/oscal"
)

func NewGemara2OSCAL(logger hclog.Logger) *cobra.Command {
	option := NewOptions()

	command := &cobra.Command{
		Use:   "gemara2oscal",
		Short: "Generate OSCAL Assessment Layer artifacts from Gemara artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := option.Complete(cmd, logger); err != nil {
				return err
			}
			return runGemara2OSCAL(cmd.Context(), option)
		},
	}

	fs := command.Flags()
	BindGemaraFlags(fs)
	fs.StringP("out", "o", "./assessment-results.json", "path to output OSCAL Assessment Results")
	fs.String(Catalog, "", "path to catalog.json")
	return command
}

func runGemara2OSCAL(ctx context.Context, option *Options) error {
	validator := validation.NewSchemaValidator()
	catalogFile, err := os.Open(option.Catalog)
	if err != nil {
		return err
	}
	defer catalogFile.Close()
	catalog, err := models.NewCatalog(catalogFile, validator)
	if err != nil {
		return fmt.Errorf("error loading catalog: %w", err)
	}

	config, err := getPolicy(option.Policy)
	if err != nil {
		return err
	}

	var plans []policy.PlanRef
	// Set loaders
	for _, catalogRef := range config.Catalogs {
		for _, plan := range catalogRef.Plans {
			// Lazily load evals
			plan.Loader = func() (*layer4.Layer4, error) {
				var l4Eval layer4.Layer4
				filePath := filepath.Clean(filepath.Join(option.EvalDir, "results", fmt.Sprintf("%s.yml", plan.Service)))
				data, err := os.ReadFile(filePath)
				if err != nil {
					return nil, err
				}

				err = yaml.Unmarshal(data, &l4Eval)
				if err != nil {
					return nil, err
				}
				return &l4Eval, nil
			}
			plans = append(plans, plan)
		}
	}

	logger := option.Logger()
	logger.Debug(fmt.Sprintf("Using catalog %s", catalog.Metadata.Title))

	gemaraProvider, err := gemara.NewGemaraValidatorFromFile(option.Policy, option.EvalDir)
	if err != nil {
		return err
	}

	inputCtx, err := gemaraProvider.Plan()
	if err != nil {
		return err
	}

	plan, err := gemerara2AssessmentPlan(plans, catalog.Metadata.Title, logger)
	if err != nil {
		return err
	}

	// Assuming the OSCAL Catalog Title would line up with the ID in the Eval
	allResults, err := gemara2AssessmentResults(plans, catalog.Metadata.Title, logger)
	if err != nil {
		return err
	}

	oscalProvider := oscal.NewOSCALValidation(plan, "REPLACE_ME", logger)
	return oscalProvider.Report(ctx, inputCtx, option.Output, allResults)
}

// TODO: Complete
func gemerara2AssessmentPlan(plans []policy.PlanRef, catalogId string, logger hclog.Logger) (*oscalTypes.AssessmentPlan, error) {
	return &oscalTypes.AssessmentPlan{}, nil
}

func gemara2AssessmentResults(plans []policy.PlanRef, catalogId string, logger hclog.Logger) ([]policy.PVPResult, error) {
	var allResults []policy.PVPResult
	var inScopePlan uint
	for _, plan := range plans {
		if plan.Plan == nil {
			logger.Debug(fmt.Sprintf("Loading plan for %s", plan.Service))
			if err := plan.Load(); err != nil {
				return nil, err
			}
		}
		if plan.Plan.CatalogID != catalogId {
			logger.Debug(fmt.Sprintf("Plan %s does not match %s. Skipping...", plan.Plan.CatalogID, catalogId))
			continue
		}
		inScopePlan++

		result := gemara.ObservationsFromEvaluation(*plan.Plan)
		allResults = append(allResults, result)
	}
	logger.Debug(fmt.Sprintf("Processed %v in scope plans", inScopePlan))
	return allResults, nil
}

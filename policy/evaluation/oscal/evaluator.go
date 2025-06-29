package oscal

import (
	"context"
	"fmt"
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/hashicorp/go-hclog"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
)

var _ evaluation.Provider = (*OSCALValidation)(nil)

type OSCALValidation struct {
	assessmentPlan *oscalTypes.AssessmentPlan
	planLoc        string
	logger         hclog.Logger
}

func NewOSCALValidation(plan *oscalTypes.AssessmentPlan, planLoc string, logger hclog.Logger) *OSCALValidation {
	return &OSCALValidation{
		assessmentPlan: plan,
		planLoc:        planLoc,
		logger:         logger,
	}
}

func NewOSCALValidationFromFile(planLoc string, logger hclog.Logger) (*OSCALValidation, error) {
	plan, err := loadPlan(planLoc)
	if err != nil {
		return nil, err
	}
	return NewOSCALValidation(plan, planLoc, logger), nil
}

func (o OSCALValidation) Report(ctx context.Context, inputCtx *evaluation.InputContext, output string, results []policy.PVPResult) error {
	assessmentResults, err := Report(ctx, inputCtx, o.planLoc, *o.assessmentPlan, results)
	if err != nil {
		return err
	}

	oscalModels := oscalTypes.OscalModels{
		AssessmentResults: assessmentResults,
	}

	// Validate before writing out
	o.logger.Info("Validating generated assessment results")
	validator := validation.NewSchemaValidator()
	if err := validator.Validate(oscalModels); err != nil {
		return err
	}

	o.logger.Info(fmt.Sprintf("Writing assessment results to %s.", output))
	err = pkg.WriteObjToJsonFile(output, oscalModels)
	if err != nil {
		return err
	}
	return nil
}

func (o OSCALValidation) Plan() (*evaluation.InputContext, error) {
	return NewContext(o.assessmentPlan)
}

func loadPlan(path string) (*oscalTypes.AssessmentPlan, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	plan, err := models.NewAssessmentPlan(file, validation.NewSchemaValidator())
	if err != nil {
		return nil, err
	}
	return plan, nil
}

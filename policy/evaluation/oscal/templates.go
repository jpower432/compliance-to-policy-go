package oscal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/hashicorp/go-hclog"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/transformers"
	"github.com/oscal-compass/oscal-sdk-go/validation"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
)

func GenerateTemplates(ctx context.Context, logger hclog.Logger, path, planLoc, applicability string) error {
	cleanedPath := filepath.Clean(path)
	data, err := os.Open(cleanedPath)
	if err != nil {
		return err
	}
	compDef, err := models.NewComponentDefinition(data, validation.NewSchemaValidator())
	if err != nil {
		return err
	}

	ap, err := transformers.ComponentDefinitionsToAssessmentPlan(ctx, []oscalTypes.ComponentDefinition{*compDef}, applicability)
	if err != nil {
		return err
	}

	oscalModels := oscalTypes.OscalModels{
		AssessmentPlan: ap,
	}
	logger.Info(fmt.Sprintf("Writing assessment plan to %s.", planLoc))
	err = pkg.WriteObjToJsonFile(planLoc, oscalModels)
	if err != nil {
		return err
	}
	return nil
}

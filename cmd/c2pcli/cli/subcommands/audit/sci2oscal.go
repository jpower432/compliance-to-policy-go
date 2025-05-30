/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package audit

import (
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

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/options"
	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/subcommands/eval"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/convert"
	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
)

func NewSCI2OSCAL(logger hclog.Logger) *cobra.Command {
	option := options.NewOptions()

	command := &cobra.Command{
		Use:   "sci2oscal",
		Short: "Generate OSCAL Assessment Layer artifacts from SCI artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := option.Complete(cmd, logger); err != nil {
				return err
			}
			return runSCI2OSCAL(option)
		},
	}

	fs := command.Flags()
	options.BindSCIFlags(fs)
	fs.StringP("out", "o", "./assessment-results.json", "path to output OSCAL Assessment Results")
	fs.String(options.Catalog, "", "path to catalog.json")
	return command
}

func runSCI2OSCAL(option *options.Options) error {
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

	policy, err := eval.GetPolicy(option.Policy)
	if err != nil {
		return err
	}
	// Set loaders
	for i := range policy.Refs {
		// Lazily load evals
		policy.Refs[i].Loader = func() (*layer4.Layer4, error) {
			var l4Eval layer4.Layer4
			filePath := filepath.Clean(filepath.Join(option.EvalDir, "results", fmt.Sprintf("%s.yml", policy.Refs[i].Service)))
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
	}
	logger := option.Logger()
	logger.Debug(fmt.Sprintf("Using catalog %s", catalog.Metadata.Title))

	// Assuming the OSCAL Catalog Title would line up with the ID in the Eval
	assessmentResults, err := convert.SCI2AssessmentResults(policy.Refs, catalog.Metadata.Title, logger)
	if err != nil {
		return err
	}

	oscalModels := oscalTypes.OscalModels{
		AssessmentResults: &assessmentResults,
	}

	// Validate before writing out
	logger.Info("Validating generated assessment results")
	if err := validator.Validate(oscalModels); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Writing assessment results to %s.", option.Output))
	err = pkg.WriteObjToJsonFile(option.Output, oscalModels)
	if err != nil {
		return err
	}

	return nil
}

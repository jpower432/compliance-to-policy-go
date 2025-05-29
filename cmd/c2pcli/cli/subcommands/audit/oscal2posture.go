/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package audit

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/options"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/posture"
)

func NewOSCAL2Posture(logger hclog.Logger) *cobra.Command {
	option := options.NewOptions()

	command := &cobra.Command{
		Use:   "oscal2posture",
		Short: "Generate Compliance Posture from OSCAL artifacts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := option.Complete(cmd, logger); err != nil {
				return err
			}
			if err := validateOSCAL2Posture(option); err != nil {
				return err
			}
			return runOSCAL2Posture(option)
		},
	}
	fs := command.Flags()
	options.BindCommonFlags(fs)
	fs.String(options.Catalog, "", "path to catalog.json")
	fs.StringP("assessment-results", "a", "./assessment-results.json", "path to assessment-results.json")
	fs.StringP(options.ComponentDefinition, "d", "", "path to component-definition.json file. This option cannot be used with --assessment-plan.")
	fs.StringP("out", "o", "-", "path to output file. Use '-' for stdout. Default '-'.")
	return command
}

// validateOSCAL2Posture runs validation specific to the OSCAL2Posture command.
func validateOSCAL2Posture(option *options.Options) error {
	var errs []error
	if option.Catalog == "" {
		errs = append(errs, &options.ConfigError{Option: options.Catalog})
	}
	if option.Definition == "" {
		errs = append(errs, &options.ConfigError{Option: options.ComponentDefinition})
	}
	return errors.Join(errs...)
}

func runOSCAL2Posture(option *options.Options) error {
	schemaValidator := validation.NewSchemaValidator()
	arFile, err := os.Open(option.AssessmentResults)
	if err != nil {
		return err
	}
	defer arFile.Close()
	assessmentResults, err := models.NewAssessmentResults(arFile, schemaValidator)
	if err != nil {
		return fmt.Errorf("error loading assessment results: %w", err)
	}

	catalogFile, err := os.Open(option.Catalog)
	if err != nil {
		return err
	}
	defer catalogFile.Close()
	catalog, err := models.NewCatalog(catalogFile, schemaValidator)
	if err != nil {
		return fmt.Errorf("error loading catalog: %w", err)
	}

	compDefFile, err := os.Open(option.Definition)
	if err != nil {
		return err
	}
	defer compDefFile.Close()
	compDef, err := models.NewComponentDefinition(compDefFile, schemaValidator)
	if err != nil {
		return fmt.Errorf("error loading component definition: %w", err)
	}

	r := posture.NewOscal2Posture(assessmentResults, catalog, compDef, option.Logger())
	data, err := r.Generate()
	if err != nil {
		return err
	}

	out := option.Output
	if out == "-" {
		fmt.Fprintln(os.Stdout, string(data))
	} else {
		return os.WriteFile(out, data, 0600)
	}
	return nil
}

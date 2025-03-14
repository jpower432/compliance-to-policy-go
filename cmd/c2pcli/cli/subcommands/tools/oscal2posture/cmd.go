/*
Copyright 2023 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oscal2posture

import (
	"fmt"
	"os"

	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
)

var logger *zap.Logger = pkg.GetLogger("oscal2posture")

func New(logger *zap.Logger) *cobra.Command {
	opts := NewOptions()

	command := &cobra.Command{
		Use:   "oscal2posture",
		Short: "Generate Compliance Posture from OSCAL artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Complete(); err != nil {
				return err
			}

			if err := opts.Validate(); err != nil {
				return err
			}
			return Run(opts, logger)
		},
	}
	opts.AddFlags(command.Flags())

	return command
}

func Run(options *Options, logger *zap.Logger) error {

	arFile, err := os.Open(options.AssessmentResults)
	if err != nil {
		return err
	}
	defer arFile.Close()
	assessmentResults, err := models.NewAssessmentResults(arFile, validation.NewSchemaValidator())
	if err != nil {
		return err
	}

	catalogFile, err := os.Open(options.Catalog)
	if err != nil {
		return err
	}
	defer catalogFile.Close()
	catalog, err := models.NewCatalog(catalogFile, validation.NewSchemaValidator())
	if err != nil {
		return err
	}

	compDefFile, err := os.Open(options.ComponentDefinition)
	if err != nil {
		return err
	}
	defer compDefFile.Close()
	compDef, err := models.NewComponentDefinition(compDefFile, validation.NewSchemaValidator())
	if err != nil {
		return err
	}

	r := framework.NewOscal2Posture(assessmentResults, catalog, compDef, logger)
	data, err := r.Generate()
	if err != nil {
		return err
	}

	if options.Out == "-" {
		fmt.Fprintln(os.Stdout, string(data))
	} else {
		return os.WriteFile(options.Out, data, os.ModePerm)
	}

	return nil
}

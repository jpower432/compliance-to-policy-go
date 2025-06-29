/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package subcommands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions_Validate(t *testing.T) {
	tests := []struct {
		name      string
		options   *Options
		wantError string
	}{
		{
			name: "Invalid/BothOptionsSet",
			options: &Options{
				OSCALOptions: OSCALOptions{
					Definition: "set",
					Plan:       "also-set",
				},
			},
			wantError: "cannot set both component-definition and assessment-plan values",
		},
		{
			name:      "Invalid/NoOptionsSet",
			options:   &Options{},
			wantError: "must set component-definition or assessment-plan",
		},
		{
			name: "Invalid/InvalidOptionsSet",
			options: &Options{
				OSCALOptions: OSCALOptions{
					Definition: "set",
				},
			},
			wantError: "\"name\" option is not set",
		},
		{
			name: "Valid/PlanSet",
			options: &Options{
				OSCALOptions: OSCALOptions{
					Plan: "also-set",
				},
			},
		},
		{
			name: "Valid/DefinitionSet",
			options: &Options{
				OSCALOptions: OSCALOptions{
					Definition: "set",
					Name:       "set",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.options.Validate()

			if test.wantError != "" {
				require.EqualError(t, err, test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}

}

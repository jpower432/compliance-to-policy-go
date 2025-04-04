/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package action

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/stretchr/testify/require"
)

func TestGetPluginIDFromComponent(t *testing.T) {
	tests := []struct {
		name      string
		component oscalTypes.DefinedComponent
		expected  string
		wantError string
	}{
		{
			name: "Valid/ExactID",
			component: oscalTypes.DefinedComponent{
				Title: "myplugin",
			},
			expected:  "myplugin",
			wantError: "",
		},
		{
			name: "Valid/WithWhiteSpace",
			component: oscalTypes.DefinedComponent{
				Title: " myplugin ",
			},
			expected:  "myplugin",
			wantError: "",
		},
		{
			name: "Valid/UpperCase",
			component: oscalTypes.DefinedComponent{
				Title: "MYPLUGIN",
			},
			expected:  "myplugin",
			wantError: "",
		},
		{
			name: "Invalid/PluginNotMatchPattern",
			component: oscalTypes.DefinedComponent{
				Title: "my plugin",
			},
			expected:  "",
			wantError: "invalid plugin id my plugin",
		},
		{
			name: "Invalid/EmptyTitle",
			component: oscalTypes.DefinedComponent{
				Title: "",
			},
			expected:  "",
			wantError: "component is missing a title",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			compAdapter := components.NewDefinedComponentAdapter(c.component)
			id, err := GetPluginIDFromComponent(compAdapter)
			if c.wantError == "" {
				require.NoError(t, err)
				require.Equal(t, c.expected, id)
			} else {
				require.EqualError(t, err, c.wantError)
			}
		})
	}
}

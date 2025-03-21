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

package kyverno

import (
	"os"
	"testing"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
	typec2pcr "github.com/oscal-compass/compliance-to-policy-go/v2/pkg/types/c2pcr"
	"github.com/stretchr/testify/assert"
)

func TestOscal2Policy(t *testing.T) {
	policyDir := pkg.PathFromPkgDirectory("./testdata/kyverno/policy-resources")
	catalogPath := pkg.PathFromPkgDirectory("./testdata/oscal/catalog.json")
	profilePath := pkg.PathFromPkgDirectory("./testdata/oscal/profile.json")
	cdPath := pkg.PathFromPkgDirectory("./testdata/kyverno/component-definition.json")
	// expectedDir := pkg.PathFromPkgDirectory("./composer/testdata/expected/c2pcr-parser-composed-policies")

	tempDirPath := pkg.PathFromPkgDirectory("./testdata/_test")
	err := os.MkdirAll(tempDirPath, os.ModePerm)
	assert.NoError(t, err, "Should not happen")
	tempDir := pkg.NewTempDirectory(tempDirPath)

	gitUtils := pkg.NewGitUtils(tempDir)

	c2pcrSpec := typec2pcr.Spec{
		Compliance: typec2pcr.Compliance{
			Name: "Test Compliance",
			Catalog: typec2pcr.ResourceRef{
				Url: catalogPath,
			},
			Profile: typec2pcr.ResourceRef{
				Url: profilePath,
			},
			ComponentDefinition: typec2pcr.ResourceRef{
				Url: cdPath,
			},
		},
		PolicyResources: typec2pcr.ResourceRef{
			Url: policyDir,
		},
		ClusterGroups: []typec2pcr.ClusterGroup{{
			Name:        "test-group",
			MatchLabels: &map[string]string{"environment": "test"},
		}},
		Binding: typec2pcr.Binding{
			Compliance:    "Test Compliance",
			ClusterGroups: []string{"test-group"},
		},
		Target: typec2pcr.Target{
			Namespace: "",
		},
	}
	c2pcrParser := NewParser(gitUtils)
	c2pcrParsed, err := c2pcrParser.Parse(c2pcrSpec)
	assert.NoError(t, err, "Should not happen")

	o2p := NewOscal2Policy(c2pcrParsed.PolicyResoureDir, tempDir)
	err = o2p.Generate(c2pcrParsed)
	assert.NoError(t, err, "Should not happen")
}

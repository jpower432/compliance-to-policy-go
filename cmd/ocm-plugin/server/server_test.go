package server

import (
	"context"
	"os"
	"testing"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
)

func TestOscal2Policy(t *testing.T) {
	policyDir := pkg.PathFromPkgDirectory("./testdata/ocm/policies")

	tempDirPath := pkg.PathFromPkgDirectory("./testdata/_test")
	err := os.MkdirAll(tempDirPath, os.ModePerm)
	assert.NoError(t, err, "Should not happen")
	tempDir := pkg.NewTempDirectory(tempDirPath)

	testPolicy := createPolicy(t)
	plugin := NewPlugin()
	plugin.config.policiesDir = policyDir
	plugin.config.namespace = "test"
	plugin.config.policySetName = "test"
	plugin.config.tempDir = tempDir.GetTempDir()
	require.NoError(t, plugin.Generate(testPolicy))

	assert.NoError(t, err, "Should not happen")
}

func TestResult2Oscal(t *testing.T) {

	policyResultsDir := pkg.PathFromPkgDirectory("./testdata/ocm/policy-results")

	tempDirPath := pkg.PathFromPkgDirectory("./testdata/_test")
	err := os.MkdirAll(tempDirPath, os.ModePerm)
	assert.NoError(t, err, "Should not happen")

	testPolicy := createPolicy(t)

	reporter := NewResultToOscal(testPolicy, policyResultsDir, "example", "example")
	_, err = reporter.GenerateResults()
	assert.NoError(t, err, "Should not happen")
}

func createPolicy(t *testing.T) []extensions.RuleSet {
	cdPath := pkg.PathFromPkgDirectory("./testdata/ocm/component-definition.json")

	file, err := os.Open(cdPath)
	require.NoError(t, err)
	defer file.Close()

	compDef, err := generators.NewComponentDefinition(file)

	require.NotNil(t, compDef)
	require.NotNil(t, compDef.Components)

	var allComponents []components.Component
	for _, comp := range *compDef.Components {
		adapter := components.NewDefinedComponentAdapter(comp)
		allComponents = append(allComponents, adapter)
	}

	store := rules.NewMemoryStore()
	require.NoError(t, store.IndexAll(allComponents))

	ruleSets, err := store.FindByComponent(context.TODO(), "OCM")
	require.NoError(t, err)
	return ruleSets
}

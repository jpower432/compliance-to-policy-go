/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"
	"errors"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/compliance-to-policy-go/v2/internal/utils"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var (
	expectedCertFileRule = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "etcd_cert_file",
			Description: "Ensure that the --cert-file argument is set as appropriate",
		},
		Checks: []extensions.Check{
			{
				ID:          "etcd_cert_file",
				Description: "Check that the --cert-file argument is set as appropriate",
			},
		},
	}
	expectedKeyFileRule = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "etcd_key_file",
			Description: "Ensure that the --key-file argument is set as appropriate",
			Parameters: []extensions.Parameter{
				{
					ID:          "file_name",
					Description: "A parameter for a file name",
				},
			},
		},
		Checks: []extensions.Check{
			{
				ID:          "etcd_key_file",
				Description: "Check that the --key-file argument is set as appropriate",
			},
		},
	}
)

func TestAggregateResults(t *testing.T) {
	inputContext := inputContextHelper(t)
	wantResults := policy.PVPResult{
		ObservationsByCheck: []policy.ObservationByCheck{
			{
				Title:       "Example",
				Description: "Example description",
				CheckID:     "test-check",
			},
		},
	}

	updatedParam := extensions.Parameter{
		ID:          "file_name",
		Description: "A parameter for a file name",
		Value:       "my_file",
	}

	updatedKeyFileRule := expectedKeyFileRule
	updatedKeyFileRule.Rule.Parameters[0] = updatedParam

	// Create pluginSet
	providerTestObj := new(policyProvider)
	providerTestObj.On("GetResults", policy.Policy{updatedKeyFileRule}).Return(wantResults, nil)
	pluginSet := map[plugin.ID]policy.Provider{
		"mypvpvalidator": providerTestObj,
	}

	testSettings := settings.NewSettings(map[string]struct{}{"etcd_key_file": {}}, map[string]string{"file_name": "my_file"})
	inputContext.Settings = testSettings

	gotResults, err := AggregateResults(context.TODO(), inputContext, pluginSet)
	require.NoError(t, err)
	providerTestObj.AssertExpectations(t)
	require.Len(t, gotResults, 1)
}

func TestAggregateResults_Multi(t *testing.T) {
	testDataPath := utils.PathFromInternalDirectory("./testdata/oscal/component-definition-heterogeneous.json")
	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := models.NewComponentDefinition(file, validation.NoopValidator{})
	require.NoError(t, err)

	var allComponents []components.Component
	for _, component := range *definition.Components {
		compAdapter := components.NewDefinedComponentAdapter(component)
		allComponents = append(allComponents, compAdapter)
	}

	inputContext, err := NewContext(allComponents)
	require.NoError(t, err)

	testSettings := settings.NewSettings(map[string]struct{}{"test_configuration_check": {}, "allowed-base-images": {}}, nil)
	inputContext.Settings = testSettings

	wantResults := policy.PVPResult{
		ObservationsByCheck: []policy.ObservationByCheck{
			{
				Title:       "Example",
				Description: "Example description",
				CheckID:     "test-check",
			},
		},
	}

	ocmRule := extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "test_configuration_check",
			Description: "Ensure deployment configuration is securely set up",
		},
		Checks: []extensions.Check{
			{
				ID: "policy-high-scan",
			},
		},
	}

	kyvernoRule := extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "allowed-base-images",
			Description: "Building images which specify a base as their origin is a good start to improving supply chain security, but over time organizations may want to build an allow list of specific base images which are allowed to be used when constructing containers. This policy ensures that a container's base, found in an OCI annotation, is in a cluster-wide allow list.",
		},
		Checks: []extensions.Check{
			{
				ID:          "allowed-base-images",
				Description: "allowed-base-images",
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		providerTestObj := new(policyProvider)
		providerTestObj.On("GetResults", policy.Policy{ocmRule}).Return(wantResults, nil)
		providerTestObj2 := new(policyProvider)
		providerTestObj2.On("GetResults", policy.Policy{kyvernoRule}).Return(wantResults, nil)

		pluginSet := map[plugin.ID]policy.Provider{
			"ocm":     providerTestObj,
			"kyverno": providerTestObj2,
		}

		gotResults, err := AggregateResults(context.TODO(), inputContext, pluginSet)
		require.NoError(t, err)
		require.Len(t, gotResults, 2)
		providerTestObj.AssertExpectations(t)
		providerTestObj2.AssertExpectations(t)
	})

	t.Run("Failing Provider", func(t *testing.T) {
		providerTestObj := new(policyProvider)
		providerTestObj.On("GetResults", policy.Policy{ocmRule}).Return(wantResults, nil)
		providerTestObj3 := new(policyProvider)
		providerTestObj3.On("GetResults", policy.Policy{kyvernoRule}).Return(policy.PVPResult{}, errors.New("failed"))

		pluginSet := map[plugin.ID]policy.Provider{
			"ocm":     providerTestObj,
			"kyverno": providerTestObj3,
		}

		gotResults, err := AggregateResults(context.Background(), inputContext, pluginSet)
		require.EqualError(t, err, "failed")

		// The number for could 0 or 1 because of randomness in the map key order.
		assert.Contains(t, []int{0, 1}, len(gotResults))

		providerTestObj.AssertExpectations(t)
		providerTestObj3.AssertExpectations(t)
	})

	t.Run("Cancellation Error", func(t *testing.T) {
		providerTestObj := new(policyProvider)
		// Mock a successful call that will be cancelled
		providerTestObj.On("GetResults", policy.Policy{ocmRule}).Return(wantResults, nil)
		providerTestObj3 := new(policyProvider)
		// This mock will have a delay, ensuring it gets cancelled
		providerTestObj3.On("GetResults", policy.Policy{kyvernoRule}).Return(policy.PVPResult{}, errors.New("should not happen"))
		providerTestObj3.delay = 500 * time.Millisecond // Assuming your mock provider has a delay field

		pluginSet := map[plugin.ID]policy.Provider{
			"ocm":     providerTestObj,
			"kyverno": providerTestObj3,
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel the context almost immediately
		cancel()

		_, err := AggregateResults(ctx, inputContext, pluginSet)

		// The error should be a context error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}

// policyProvider is a mocked implementation of policy.Provider.
type policyProvider struct {
	mock.Mock
	delay time.Duration
}

func (p *policyProvider) Configure(_ context.Context, option map[string]string) error {
	args := p.Called(option)
	return args.Error(0)
}

func (p *policyProvider) Generate(_ context.Context, policyRules policy.Policy) error {
	sort.SliceStable(policyRules, func(i, j int) bool {
		return policyRules[i].Rule.ID > policyRules[j].Rule.ID
	})
	args := p.Called(policyRules)
	return args.Error(0)
}

func (p *policyProvider) GetResults(ctx context.Context, policyRules policy.Policy) (policy.PVPResult, error) {
	sort.SliceStable(policyRules, func(i, j int) bool {
		return policyRules[i].Rule.ID > policyRules[j].Rule.ID
	})
	args := p.Called(policyRules)

	select {
	case <-ctx.Done():
		return policy.PVPResult{}, ctx.Err()
	case <-time.After(p.delay): // Simulate completing work
		return args.Get(0).(policy.PVPResult), args.Error(1)
	}
}

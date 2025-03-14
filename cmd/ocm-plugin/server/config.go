/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package server

type Config struct {
	policiesDir      string            `mapstructure:"policy-dir"`
	policyResultsDir string            `mapstructure:"polciy-results-dir"`
	tempDir          string            `mapstructure:"temp-dir"`
	outputDir        string            `mapstructure:"output-dir"`
	namespace        string            `mapstructure:"namespace"`
	policySetName    string            `mapstructure:"policy-set-name"`
	clusterSelectors map[string]string `mapstructure:"cluster-selectors"`
}

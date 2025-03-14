/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package server

type Config struct {
	PoliciesDir      string            `mapstructure:"policy-dir"`
	PolicyResultsDir string            `mapstructure:"polciy-results-dir"`
	TempDir          string            `mapstructure:"temp-dir"`
	OutputDir        string            `mapstructure:"output-dir"`
	Namespace        string            `mapstructure:"Namespace"`
	PolicySetName    string            `mapstructure:"policy-set-name"`
	clusterSelectors map[string]string `mapstructure:"cluster-selectors"`
}

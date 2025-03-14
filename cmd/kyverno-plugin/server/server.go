/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var _ policy.Provider = (*Plugin)(nil)

var logger *zap.Logger = pkg.GetLogger("kyverno")

type Plugin struct {
	PoliciesDir      string `mapstructure:"policy-dir"`
	PolicyResultsDir string `mapstructure:"policy-results-dir"`
	TempDir          string `mapstructure:"temp-dir"`
	OutputDir        string `mapstructure:"output-dir"`
	logger           *zap.Logger
}

func NewPlugin() *Plugin {
	return &Plugin{
		logger: logger,
	}
}

func (p *Plugin) Configure(m map[string]string) error {
	return mapstructure.Decode(m, &p)
}

func (p *Plugin) Generate(pl policy.Policy) error {
	fmt.Println(p.PoliciesDir)
	tmpdir := pkg.NewTempDirectory(p.TempDir)
	composer := NewOscal2Policy(p.PoliciesDir, tmpdir)
	if err := composer.Generate(pl); err != nil {
		return err
	}

	if p.OutputDir != "" {
		if err := composer.CopyAllTo(p.OutputDir); err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) GetResults(pl policy.Policy) (policy.PVPResult, error) {
	results := NewResultToOscal(pl, p.PolicyResultsDir, p.logger)
	return results.GenerateResults()
}

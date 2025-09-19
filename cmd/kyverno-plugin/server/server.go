/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/complytime/complybeacon/proofwatch"
	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/otel"

	"github.com/oscal-compass/compliance-to-policy-go/v2/internal/utils"
	"github.com/oscal-compass/compliance-to-policy-go/v2/logging"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var (
	_      policy.Provider = (*Plugin)(nil)
	logger hclog.Logger    = logging.NewPluginLogger()
	tracer                 = otel.Tracer("kyverno")
	meter                  = otel.Meter("kyverno")
)

func Logger() hclog.Logger {
	return logger
}

type Plugin struct {
	config  Config
	watcher *proofwatch.ProofWatch
}

func NewPlugin() (*Plugin, error) {
	watcher, err := proofwatch.NewProofWatch("kyverno", meter)
	if err != nil {
		return nil, err
	}
	return &Plugin{watcher: watcher}, nil
}

func (p *Plugin) Configure(ctx context.Context, m map[string]string) error {
	ctx, span := tracer.Start(ctx, "plugin.Configure")
	defer span.End()
	if err := mapstructure.Decode(m, &p.config); err != nil {
		return errors.New("error decoding configuration")
	}
	return p.config.Validate()
}

func (p *Plugin) Generate(ctx context.Context, pl policy.Policy) error {
	ctx, span := tracer.Start(ctx, "plugin.Generate")
	defer span.End()

	logger.Debug(fmt.Sprintf("Using resources from %s", p.config.PoliciesDir))
	tmpdir := utils.NewTempDirectory(p.config.TempDir)
	composer := NewOscal2Policy(p.config.PoliciesDir, tmpdir)
	if err := composer.Generate(pl); err != nil {
		return err
	}

	if p.config.OutputDir != "" {
		if err := composer.CopyAllTo(p.config.OutputDir); err != nil {
			return err
		}
		logger.Debug(fmt.Sprintf("Copied outputs to %s", p.config.OutputDir))
	}
	return nil
}

func (p *Plugin) GetResults(ctx context.Context, pl policy.Policy) (policy.PVPResult, error) {
	ctx, span := tracer.Start(ctx, "plugin.GetResults")
	defer span.End()
	results := NewResultToOscal(pl, p.config.PolicyResultsDir)
	return results.GenerateResults(ctx, p.watcher)
}

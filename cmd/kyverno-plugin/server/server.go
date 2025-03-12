package server

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var _ policy.Provider = (*Plugin)(nil)

var logger *zap.Logger = pkg.GetLogger("kyverno")

type Plugin struct {
	policiesDir      string `mapstructure:"policy-dir"`
	policyResultsDir string `mapstructure:"polciy-results-dir"`
	tempDir          string `mapstructure:"temp-dir"`
	outputDir        string `mapstructure:"output-dir"`
	logger           *zap.Logger
}

func NewPlugin() *Plugin {
	return &Plugin{
		logger: logger,
	}
}

func (p *Plugin) Configure(m map[string]string) error {
	for k, v := range m {
		viper.Set(k, v)
	}
	return viper.Unmarshal(p)
}

func (p *Plugin) Generate(pl policy.Policy) error {
	tmpdir := pkg.NewTempDirectory(p.tempDir)
	composer := NewOscal2Policy(p.policiesDir, tmpdir)
	if err := composer.Generate(pl); err != nil {
		return err
	}

	if p.outputDir != "" {
		if err := composer.CopyAllTo(p.outputDir); err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) GetResults(pl policy.Policy) (policy.PVPResult, error) {
	results := NewResultToOscal(pl, p.policyResultsDir, p.logger)
	return results.GenerateResults()
}

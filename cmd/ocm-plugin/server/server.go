package server

import (
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var _ policy.Provider = (*Plugin)(nil)

type Plugin struct {
	config             Config
	policyGeneratorDir string
}

func NewPlugin() *Plugin {
	return &Plugin{
		config: Config{},
	}
}

func (p *Plugin) Configure(m map[string]string) error {
	for k, v := range m {
		viper.Set(k, v)
	}
	return viper.Unmarshal(&p.config)
}

func (p *Plugin) Generate(pl policy.Policy) error {
	tmpdir := pkg.NewTempDirectory(p.config.tempDir)
	composer := NewComposerByTempDirectory(p.config.policyResultsDir, tmpdir)
	if err := composer.ComposeByPolicies(pl, p.config); err != nil {
		panic(err)
	}
	policySet, err := composer.GeneratePolicySet()
	if err != nil {
		panic(err)
	}

	for _, resource := range (*policySet).Resources() {
		name := resource.GetName()
		kind := resource.GetKind()
		namespace := resource.GetNamespace()
		yamlByte, err := resource.AsYAML()
		if err != nil {
			panic(err)
		}
		fnamesTokens := []string{kind, namespace, name}
		fname := strings.Join(fnamesTokens, ".") + ".yaml"
		if err := os.WriteFile(p.config.outputDir+"/"+fname, yamlByte, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if p.policyGeneratorDir != "" {
		if err := composer.CopyAllTo(p.policyGeneratorDir); err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) GetResults(pl policy.Policy) (policy.PVPResult, error) {
	results := NewResultToOscal(pl, p.config.policyResultsDir, p.config.namespace, p.config.policySetName)
	return results.GenerateResults()
}

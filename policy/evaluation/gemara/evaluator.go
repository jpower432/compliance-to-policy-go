package gemara

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/revanite-io/sci/layer2"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy/evaluation"
)

var _ evaluation.Provider = (*GemaraValidation)(nil)

type GemaraValidation struct {
	policyConfig policy.Config
	evalDir      string
}

func NewGemaraValidator(config policy.Config, evalDir string) *GemaraValidation {
	return &GemaraValidation{
		policyConfig: config,
		evalDir:      evalDir,
	}
}

func NewGemaraValidatorFromFile(path, evalDir string) (*GemaraValidation, error) {
	policyConfig, err := getPolicy(path)
	if err != nil {
		return nil, err
	}
	for i := range policyConfig.Catalogs {

		catalog, err := getCatalog(policyConfig.Catalogs[i].CatalogID)
		if err != nil {
			return nil, err
		}

		// Lazily load
		for j := range policyConfig.Catalogs[i].Plans {
			filePath := filepath.Clean(filepath.Join(evalDir, fmt.Sprintf("%s-%s.yml", policyConfig.Catalogs[i].Plans[j].PluginID, catalog.Metadata.Id)))
			policyConfig.Catalogs[i].Plans[j].Loader = func() (*layer4.Layer4, error) {
				var l4Eval layer4.Layer4
				file, err := os.Open(filePath)
				if err != nil {
					return nil, err
				}
				decoder := yaml.NewDecoder(file)
				err = decoder.Decode(&l4Eval)
				if err != nil {
					return nil, err
				}
				return &l4Eval, nil
			}
		}
	}
	return NewGemaraValidator(policyConfig, evalDir), nil
}

func (l *GemaraValidation) Plan() (*evaluation.InputContext, error) {
	var plans []policy.PlanRef
	for _, catalog := range l.policyConfig.Catalogs {
		for _, plan := range catalog.Plans {
			plans = append(plans, plan)
		}
	}
	return NewContextFromPlanRefs(plans...)
}

func (l *GemaraValidation) Report(ctx context.Context, inputCtx *evaluation.InputContext, output string, results []policy.PVPResult) error {
	for _, catalogRef := range l.policyConfig.Catalogs {
		catalog, err := getCatalog(catalogRef.CatalogID)
		if err != nil {
			return err
		}
		catalogRef.Catalog = &catalog
		plans, err := Layer4FromResults(ctx, inputCtx, results, &catalogRef)
		if err != nil {
			return err
		}
		for _, ref := range plans {
			data, err := MarshalConfigWithFunction(ref.Plan)
			if err != nil {
				return err
			}

			// Write the resulting evaluation for each service to a new file
			filePath := filepath.Clean(filepath.Join(output, fmt.Sprintf("%s.yml", ref.Service)))
			if err := os.WriteFile(filePath, data, os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}

func getCatalog(filepath string) (layer2.Layer2, error) {
	var catalog layer2.Layer2
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return catalog, err
	}
	err = yaml.Unmarshal(yamlFile, &catalog)
	if err != nil {
		return catalog, err
	}
	return catalog, nil
}

func getPolicy(filepath string) (policy.Config, error) {
	var p policy.Config
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

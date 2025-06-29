package gemara

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/revanite-io/sci/layer2"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

func GenerateTemplates(input, output string) error {
	policyConfig, err := getPolicy(input)
	if err != nil {
		return err
	}
	_, err = generateTemplates(&policyConfig, output)
	return err
}

func generateTemplates(getPolicy *policy.Config, evalDir string) (foundRefs []policy.PlanRef, err error) {
	// FIXME: Duplicates need to be removed one per validator and Layer 2 catalog
	for i := range getPolicy.Catalogs {
		catalog, err := getCatalog(getPolicy.Catalogs[i].CatalogID)
		if err != nil {
			return foundRefs, err
		}
		getPolicy.Catalogs[i].Catalog = &catalog

		// Set loaders
		for _, ref := range getPolicy.Catalogs[i].Plans {
			// Config are under the plugin name <plugin-id>-<catalog-id>.yml
			filePath := filepath.Clean(filepath.Join(evalDir, fmt.Sprintf("%s-%s.yml", ref.PluginID, catalog.Metadata.Id)))
			if _, err := os.Stat(filePath); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					if err := generateNewEvalTemplate(catalog, filePath); err != nil {
						return foundRefs, err
					}
					// Do not generate policy if creating a template from scratch
					continue
				}
				return foundRefs, err
			}
			foundRefs = append(foundRefs, ref)
		}
	}
	return foundRefs, nil
}

func generateNewEvalTemplate(catalog layer2.Layer2, filePath string) error {
	eval := layer4.NewEvaluation(catalog)
	data, err := yaml.Marshal(eval)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, os.ModePerm)
}

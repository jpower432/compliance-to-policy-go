package mapping

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/rules"
	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// Mapping is a helper function for plugins who need to map rules to checks from
// mapping configuration
type Mapping interface {
	MethodsForRequirement(requirementId string) ([]string, error)
	CompletePolicy(policy *policy.Policy) error
}

var (
	_ Mapping = (*OSCALValidation)(nil)
	_ Mapping = (*Layer4Eval)(nil)
)

type OSCALValidation struct {
	store rules.Store
}

func NewOSCALValidation(store rules.Store) OSCALValidation {
	return OSCALValidation{
		store: store,
	}
}

func NewOSCALValidationFromFile(path string) (OSCALValidation, error) {
	cleanedPath := filepath.Clean(path)
	data, err := os.Open(cleanedPath)
	if err != nil {
		return OSCALValidation{}, err
	}
	compdef, err := models.NewComponentDefinition(data, validation.NewSchemaValidator())
	if err != nil {
		return OSCALValidation{}, err
	}

	memoryStore := rules.NewMemoryStore()

	var allComponents []components.Component
	if compdef.Components == nil {
		return OSCALValidation{}, errors.New("no components in component definition")
	}
	for _, component := range *compdef.Components {
		compAdapter := components.NewDefinedComponentAdapter(component)
		allComponents = append(allComponents, compAdapter)
	}

	if err := memoryStore.IndexAll(allComponents); err != nil {
		return OSCALValidation{}, err
	}

	return NewOSCALValidation(memoryStore), nil
}

func (o OSCALValidation) CompletePolicy(policy *policy.Policy) error {
	//TODO implement me
	panic("implement me")
}

func (o OSCALValidation) MethodsForRequirement(requirementId string) ([]string, error) {
	rule, err := o.store.GetByRuleID(context.Background(), requirementId)
	if err != nil {
		return nil, err
	}
	var methods []string
	for _, check := range rule.Checks {
		methods = append(methods, check.ID)
	}
	return methods, nil
}

type Layer4Eval struct {
	eval layer4.Layer4
}

func NewLayer4Eval(eval layer4.Layer4) Layer4Eval {
	return Layer4Eval{
		eval: eval,
	}
}

func NewLayer4EvalFromFile(path string) (layer4.Layer4, error) {
	var l4Eval layer4.Layer4
	cleanedPath := filepath.Clean(path)
	data, err := os.ReadFile(cleanedPath)
	if err != nil {
		return l4Eval, err
	}
	err = yaml.Unmarshal(data, &l4Eval)
	if err != nil {
		return l4Eval, err
	}
	return l4Eval, nil
}

func (l Layer4Eval) CompletePolicy(policy *policy.Policy) error {
	//TODO implement me
	panic("implement me")
}

func (l Layer4Eval) MethodsForRequirement(requirementId string) ([]string, error) {
	var methods []string
	for _, eval := range l.eval.ControlEvaluations {
		for _, assessment := range eval.Assessments {
			if assessment.RequirementID == requirementId {
				for _, method := range assessment.Methods {
					methods = append(methods, method.Name)
				}
				break
			}
		}
	}
	return methods, nil
}

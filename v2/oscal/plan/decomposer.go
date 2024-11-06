package plan

// Decomposer reads and process information in an Assessment Plan based on input

import (
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/rules"
)

type Decomposer struct {
	plan oscalTypes.AssessmentPlan
}

func NewDecomposer(plan oscalTypes.AssessmentPlan) *Decomposer {
	return &Decomposer{
		plan: plan,
	}
}

func FromComponents(components []oscalTypes.DefinedComponent) *Decomposer {
	// Components to Assessment Plan Conversion
	return &Decomposer{}
}

func (d *Decomposer) Location() string {
	return ""
}

func (d *Decomposer) SetControlSource() {

}

func (d *Decomposer) SetComponents() {

}

// GetValidationComponents gets all validation component ids
func (d *Decomposer) GetValidationComponents() (ids []string, err error) {
	return nil, nil
}

// GetRuleData gets the rule data for each validation component
func (d *Decomposer) GetRuleData(validationID string) []rules.RuleSet {
	return nil
}

func (d *Decomposer) RuleSetByCheck(checkId string) (rules.RuleSet, error) {
	return rules.RuleSet{}, nil
}

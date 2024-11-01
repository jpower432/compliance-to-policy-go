package plan

// TODO: jpower432: Perhaps explore functional options here?

// Decomposer reads and process information in an Assessment Plan based on input
type Decomposer struct {
}

func NewDecomposer() *Decomposer {
	return &Decomposer{}
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
func (d *Decomposer) GetRuleData(validationID string) {

}

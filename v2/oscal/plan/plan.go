package plan

type Plan struct {
	Location string
	RuleFinder
}

func NewPlan(location string, store RuleFinder) *Plan {
	return &Plan{
		Location:   location,
		RuleFinder: store,
	}
}

func (p *Plan) GetPolicy() (Policy, error) {

	rulSets, err := p.All()
	if err != nil {
		return Policy{}, err
	}

	parameters := make([]*Parameter, 0)
	for _, rule := range rulSets {
		parameters = append(parameters, rule.Parameter)
	}
	policy := Policy{rulSets, parameters}
	return policy, nil
}

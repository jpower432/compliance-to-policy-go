package policy

import (
	"encoding/json"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/observations"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/rules"
)

// Provider defines methods for a policy validation engine
type Provider interface {
	GetSchema() ([]byte, error)
	UpdateConfiguration(message json.RawMessage) error
	Generate(policy rules.Policy) error
	GetResults() (observations.PVPResult, error)
}

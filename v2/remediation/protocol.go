package remediation

import (
	"encoding/json"
)

// Provider defines methods for a remediation engine
type Provider interface {
	GetSchema() ([]byte, error)
	UpdateConfiguration(message json.RawMessage) error
	Generate() error
	Remediate() error
}

package providers

import (
	"encoding/json"
)

// RemediationProvider defines methods for a remediation engine
type RemediationProvider interface {
	GetSchema() ([]byte, error)
	UpdateConfiguration(message json.RawMessage) error
	Generate() error
	Remediate() error
}

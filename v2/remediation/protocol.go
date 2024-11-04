package remediation

import (
	"encoding/json"
)

// Engine defines methods for a remediation engine
type Engine interface {
	GetSchema() ([]byte, error)
	UpdateConfiguration(message json.RawMessage) error
	Generate() error
	Remediate() error
}

package policy

import (
	"encoding/json"
)

// Engine defines methods for a policy validation engine
type Engine interface {
	GetSchema() ([]byte, error)
	UpdateConfiguration(message json.RawMessage) error
	Generate(policy Policy) error
	GetResults() (PVPResult, error)
}

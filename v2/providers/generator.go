package providers

import (
	"encoding/json"
)

/*
GenerationProvider defines methods for a policy generation engine.
Use cases include:
 1. Generation policy of artifacts using information from OSCAL objects
*/
type GenerationProvider interface {
	GetSchema() ([]byte, error)
	UpdateConfiguration(message json.RawMessage) error
	Generate(policy Policy) error
}

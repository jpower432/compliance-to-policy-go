package providers

import (
	"encoding/json"

	. "github.com/oscal-compass/compliance-to-policy-go/v2/oscal"
)

/*
		PolicyProvider defines methods for a policy validation engine.
		Defined uses cases include the following:
			1. A scanning plugin may contact a remote API for scanning
			2. A scanning plugin may exec out to another tool for scanning in a new process
	        3. A scanning plugin may be a self-contained scanning tool
*/
type PolicyProvider interface {
	GetSchema() ([]byte, error)
	UpdateConfiguration(message json.RawMessage) error
	GetResults() (PVPResult, error)
	//GetResults(policy Policy) (PVPResult, error)
}

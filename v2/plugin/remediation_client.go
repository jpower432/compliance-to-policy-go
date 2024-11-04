package plugin

import (
	"encoding/json"
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/remediation"
)

// Client must return an implementation of the corresponding interface that communicates
// over an RPC client.
var _ remediation.Engine = &remediationClient{}

type remediationClient struct {
	client proto.RemediationEngineClient
}

func (r remediationClient) GetSchema() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r remediationClient) UpdateConfiguration(message json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}

func (r remediationClient) Generate() error {
	//TODO implement me
	panic("implement me")
}

func (r remediationClient) Remediate() error {
	//TODO implement me
	panic("implement me")
}

package plugin

import (
	"context"
	"encoding/json"
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/observations"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// Client must return an implementation of the corresponding interface that communicates
// over an RPC client.
var _ policy.Provider = &pvpClient{}

type pvpClient struct {
	client proto.PolicyEngineClient
}

func (p *pvpClient) GetSchema() ([]byte, error) {
	resp, err := p.client.GetSchema(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.JsonSchema, err
}

func (p *pvpClient) UpdateConfiguration(message json.RawMessage) error {
	_, err := p.client.UpdateConfiguration(context.Background(), &proto.ConfigureRequest{
		Config: message,
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *pvpClient) Generate(policy plan.Policy) error {
	request := &proto.GenerateRequest{
		Policy: policy.ToProto(),
	}
	_, err := p.client.Generate(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

func (p *pvpClient) GetResults() (observations.PVPResult, error) {
	resp, err := p.client.GetResults(context.Background(), &proto.Empty{})
	if err != nil {
		return observations.PVPResult{}, err
	}
	pvpResult := observations.FromProto(resp.Result)
	return pvpResult, nil
}

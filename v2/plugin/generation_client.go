package plugin

import (
	"context"
	"encoding/json"

	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal"
	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

// Client must return an implementation of the corresponding interface that communicates
// over an RPC client.
var _ providers.GenerationProvider = (*generationClient)(nil)

type generationClient struct {
	client proto.GenerationEngineClient
}

func (g *generationClient) GetSchema() ([]byte, error) {
	resp, err := g.client.GetSchema(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.JsonSchema, err
}

func (g *generationClient) UpdateConfiguration(message json.RawMessage) error {
	_, err := g.client.UpdateConfiguration(context.Background(), &proto.ConfigureRequest{
		Config: message,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *generationClient) Generate(policy oscal.Policy) error {
	request := &proto.GenerateRequest{
		Policy: policy.ToProto(),
	}
	_, err := g.client.Generate(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

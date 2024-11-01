package plugin

import (
	"context"
	proto "github.com/oscal-compass/compliance-to-policy/v2/api/proto/v1alpha1"
)

var _ proto.PolicyEngineServer = &pvpService{}
var _ PolicyEngine = &pvpClient{}

type pvpClient struct {
	client proto.PolicyEngineClient
}

func (m *pvpClient) GetSchema() error {
	//TODO implement me
	panic("implement me")
}

func (m *pvpClient) UpdateConfiguration() error {
	//TODO implement me
	panic("implement me")
}

func (m *pvpClient) Generate() error {
	_, err := m.client.Generate(context.Background(), &proto.GenerateRequest{})
	return err
}

func (m *pvpClient) GetResults() error {
	_, err := m.client.GetResults(context.Background(), &proto.Empty{})
	if err != nil {
		return err
	}
	return nil
}

type pvpService struct {
	proto.UnimplementedPolicyEngineServer
	Impl PolicyEngine
}

func FromPVP(pe PolicyEngine) proto.PolicyEngineServer {
	return &pvpService{
		Impl: pe,
	}
}

func (m *pvpService) GetSchema(ctx context.Context, empty *proto.Empty) (*proto.GetSchemaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *pvpService) UpdateConfiguration(ctx context.Context, request *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *pvpService) Generate(
	ctx context.Context,
	req *proto.GenerateRequest) (*proto.GenerateResponse, error) {

	err := m.Impl.Generate()
	if err != nil {
		return nil, err
	}
	return &proto.GenerateResponse{}, nil
}

func (m *pvpService) GetResults(
	ctx context.Context,
	req *proto.Empty) (*proto.ResultsResponse, error) {
	err := m.Impl.GetResults()
	if err != nil {
		return nil, err
	}
	return &proto.ResultsResponse{}, nil
}

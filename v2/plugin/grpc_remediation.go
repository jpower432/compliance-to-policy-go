package plugin

import (
	"context"
	proto "github.com/oscal-compass/compliance-to-policy/v2/api/proto/v1alpha1"
)

var _ proto.RemediationEngineServer = &remediationService{}

var _ RemediationEngine = &remediationClient{}

type remediationClient struct {
	client proto.RemediationEngineClient
}

func (m *remediationClient) GetSchema() error {
	//TODO implement me
	panic("implement me")
}

func (m *remediationClient) UpdateConfiguration() error {
	//TODO implement me
	panic("implement me")
}

func (m *remediationClient) Generate() error {
	_, err := m.client.Generate(context.Background(), &proto.GenerateRequest{})
	return err
}

func (m *remediationClient) Remediate() error {
	_, err := m.client.Remediate(context.Background(), &proto.RemediationRequest{})
	if err != nil {
		return err
	}
	return nil
}

func FromRemediation(rm RemediationEngine) proto.RemediationEngineServer {
	return &remediationService{
		Impl: rm,
	}
}

type remediationService struct {
	proto.RemediationEngineServer
	Impl RemediationEngine
}

func (m *remediationService) GetSchema(ctx context.Context, empty *proto.Empty) (*proto.GetSchemaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *remediationService) UpdateConfiguration(ctx context.Context, request *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *remediationService) Generate(
	ctx context.Context,
	req *proto.GenerateRequest) (*proto.GenerateResponse, error) {

	err := m.Impl.Generate()
	if err != nil {
		return nil, err
	}
	return &proto.GenerateResponse{}, nil
}

func (m *remediationService) Remediate(
	ctx context.Context,
	req *proto.RemediationRequest) (*proto.RemediationResponse, error) {
	err := m.Impl.Remediate()
	if err != nil {
		return nil, err
	}
	return &proto.RemediationResponse{}, nil
}

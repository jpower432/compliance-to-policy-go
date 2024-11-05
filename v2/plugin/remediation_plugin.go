package plugin

import (
	"context"
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/remediation"
)

// Plugin must return an RPC server for this plugin type.
var _ proto.RemediationEngineServer = &remediationService{}

type remediationService struct {
	proto.UnimplementedRemediationEngineServer
	Impl remediation.Provider
}

func FromRemediation(rm remediation.Provider) proto.RemediationEngineServer {
	return &remediationService{
		Impl: rm,
	}
}

func (r remediationService) GetSchema(ctx context.Context, empty *proto.Empty) (*proto.GetSchemaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r remediationService) UpdateConfiguration(ctx context.Context, request *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r remediationService) Generate(ctx context.Context, request *proto.GenerateRequest) (*proto.GenerateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r remediationService) Remediate(ctx context.Context, request *proto.RemediationRequest) (*proto.RemediationResponse, error) {
	//TODO implement me
	panic("implement me")
}

package plugin

import (
	"context"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"

	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal/plan"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// Plugin must return an RPC server for this plugin type.
var _ proto.PolicyEngineServer = &pvpService{}

type pvpService struct {
	proto.UnimplementedPolicyEngineServer
	Impl policy.Provider
}

func FromPVP(pe policy.Provider) proto.PolicyEngineServer {
	return &pvpService{
		Impl: pe,
	}
}

func (p *pvpService) GetSchema(ctx context.Context, empty *proto.Empty) (*proto.GetSchemaResponse, error) {
	schema, err := p.Impl.GetSchema()
	if err != nil {
		return &proto.GetSchemaResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp := &proto.GetSchemaResponse{
		JsonSchema: schema,
	}
	return resp, nil
}

func (p *pvpService) UpdateConfiguration(ctx context.Context, request *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	if err := p.Impl.UpdateConfiguration(request.GetConfig()); err != nil {
		return &proto.ConfigureResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.ConfigureResponse{}, nil
}

func (p *pvpService) Generate(
	ctx context.Context,
	request *proto.GenerateRequest) (*proto.GenerateResponse, error) {

	policies := plan.FromProto(request.GetPolicy())
	if err := p.Impl.Generate(policies); err != nil {
		return &proto.GenerateResponse{Error: "I have errored"}, status.Error(codes.Internal, err.Error())
	}

	return &proto.GenerateResponse{Error: "This is a test"}, nil
}

func (p *pvpService) GetResults(
	ctx context.Context,
	request *proto.Empty) (*proto.ResultsResponse, error) {
	result, err := p.Impl.GetResults()
	if err != nil {
		return &proto.ResultsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.ResultsResponse{Result: result.ToProto()}, nil
}

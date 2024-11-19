package plugin

import (
	"context"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"

	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/oscal"
	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

// Plugin must return an RPC server for this plugin type.
var _ proto.GenerationEngineServer = (*generationService)(nil)

type generationService struct {
	proto.UnimplementedGenerationEngineServer
	Impl providers.GenerationProvider
}

func FromGenerator(gs providers.GenerationProvider) proto.GenerationEngineServer {
	return &generationService{
		Impl: gs,
	}
}

func (g *generationService) GetSchema(ctx context.Context, empty *proto.Empty) (*proto.GetSchemaResponse, error) {
	schema, err := g.Impl.GetSchema()
	if err != nil {
		return &proto.GetSchemaResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp := &proto.GetSchemaResponse{
		JsonSchema: schema,
	}
	return resp, nil
}

func (g *generationService) UpdateConfiguration(ctx context.Context, request *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	if err := g.Impl.UpdateConfiguration(request.GetConfig()); err != nil {
		return &proto.ConfigureResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.ConfigureResponse{}, nil
}

func (g *generationService) Generate(
	ctx context.Context,
	request *proto.GenerateRequest) (*proto.GenerateResponse, error) {

	policies := oscal.NewPolicyFromProto(request.GetPolicy())
	if err := g.Impl.Generate(policies); err != nil {
		return &proto.GenerateResponse{Error: "I have errored"}, status.Error(codes.Internal, err.Error())
	}

	return &proto.GenerateResponse{Error: "This is a test"}, nil
}

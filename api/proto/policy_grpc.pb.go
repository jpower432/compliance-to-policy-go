// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: api/proto/policy.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PolicyEngine_Genererate_FullMethodName = "/protocols.PolicyEngine/Genererate"
	PolicyEngine_GetResults_FullMethodName = "/protocols.PolicyEngine/GetResults"
)

// PolicyEngineClient is the client API for PolicyEngine service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// get policy results from PVP
type PolicyEngineClient interface {
	Genererate(ctx context.Context, in *PolicyRequest, opts ...grpc.CallOption) (*GenerateResponse, error)
	GetResults(ctx context.Context, in *PolicyRequest, opts ...grpc.CallOption) (*ResultsResponse, error)
}

type policyEngineClient struct {
	cc grpc.ClientConnInterface
}

func NewPolicyEngineClient(cc grpc.ClientConnInterface) PolicyEngineClient {
	return &policyEngineClient{cc}
}

func (c *policyEngineClient) Genererate(ctx context.Context, in *PolicyRequest, opts ...grpc.CallOption) (*GenerateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateResponse)
	err := c.cc.Invoke(ctx, PolicyEngine_Genererate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyEngineClient) GetResults(ctx context.Context, in *PolicyRequest, opts ...grpc.CallOption) (*ResultsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResultsResponse)
	err := c.cc.Invoke(ctx, PolicyEngine_GetResults_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PolicyEngineServer is the server API for PolicyEngine service.
// All implementations must embed UnimplementedPolicyEngineServer
// for forward compatibility.
//
// get policy results from PVP
type PolicyEngineServer interface {
	Genererate(context.Context, *PolicyRequest) (*GenerateResponse, error)
	GetResults(context.Context, *PolicyRequest) (*ResultsResponse, error)
	mustEmbedUnimplementedPolicyEngineServer()
}

// UnimplementedPolicyEngineServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPolicyEngineServer struct{}

func (UnimplementedPolicyEngineServer) Genererate(context.Context, *PolicyRequest) (*GenerateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Genererate not implemented")
}
func (UnimplementedPolicyEngineServer) GetResults(context.Context, *PolicyRequest) (*ResultsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResults not implemented")
}
func (UnimplementedPolicyEngineServer) mustEmbedUnimplementedPolicyEngineServer() {}
func (UnimplementedPolicyEngineServer) testEmbeddedByValue()                      {}

// UnsafePolicyEngineServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PolicyEngineServer will
// result in compilation errors.
type UnsafePolicyEngineServer interface {
	mustEmbedUnimplementedPolicyEngineServer()
}

func RegisterPolicyEngineServer(s grpc.ServiceRegistrar, srv PolicyEngineServer) {
	// If the following call pancis, it indicates UnimplementedPolicyEngineServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PolicyEngine_ServiceDesc, srv)
}

func _PolicyEngine_Genererate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyEngineServer).Genererate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyEngine_Genererate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyEngineServer).Genererate(ctx, req.(*PolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyEngine_GetResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyEngineServer).GetResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyEngine_GetResults_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyEngineServer).GetResults(ctx, req.(*PolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PolicyEngine_ServiceDesc is the grpc.ServiceDesc for PolicyEngine service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PolicyEngine_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocols.PolicyEngine",
	HandlerType: (*PolicyEngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Genererate",
			Handler:    _PolicyEngine_Genererate_Handler,
		},
		{
			MethodName: "GetResults",
			Handler:    _PolicyEngine_GetResults_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/policy.proto",
}
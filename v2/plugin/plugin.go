package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/providers"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "C2P_PLUGIN",
	MagicCookieValue: "0d28d48f-36e6-4026-ab12-6eae611b803b",
}

const (
	// PVPPluginName is used to dispense a policy plugin type
	PVPPluginName = "pvp"
	// RemediationPluginName is used to dispense a remediation plugin type
	RemediationPluginName = "remediation"
	// GenerationPluginName is used to dispense a generation plugin type
	GenerationPluginName = "generation"
)

// SupportedPlugins is the map of plugins we can dispense.
var SupportedPlugins = map[string]plugin.Plugin{
	PVPPluginName:         &PVPPlugin{},
	RemediationPluginName: &RemediationPlugin{},
	GenerationPluginName:  &GeneratorPlugin{},
}

var _ plugin.GRPCPlugin = &PVPPlugin{}
var _ plugin.GRPCPlugin = &RemediationPlugin{}

// Below types are only used for plugins that are written in Go.

// PVPPlugin is concrete implementation of the policy Provider written in Go.
type PVPPlugin struct {
	plugin.Plugin
	Impl providers.PolicyProvider
}

func (p *PVPPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPolicyEngineServer(s, FromPVP(p.Impl))
	return nil
}

func (p *PVPPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &pvpClient{client: proto.NewPolicyEngineClient(c)}, nil
}

// RemediationPlugin is concrete implementation of the remediation Provider written in Go.
type RemediationPlugin struct {
	plugin.Plugin
	Impl providers.RemediationProvider
}

func (p *RemediationPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRemediationEngineServer(s, FromRemediation(p.Impl))
	return nil
}

func (p *RemediationPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &remediationClient{client: proto.NewRemediationEngineClient(c)}, nil
}

// GeneratorPlugin is concrete implementation of the generation Provider written in Go.
type GeneratorPlugin struct {
	plugin.Plugin
	Impl providers.GenerationProvider
}

func (p *GeneratorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterGenerationEngineServer(s, FromGenerator(p.Impl))
	return nil
}

func (p *GeneratorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &generationClient{client: proto.NewGenerationEngineClient(c)}, nil
}

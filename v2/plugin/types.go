package plugin

import (
	"context"
	"github.com/hashicorp/go-plugin"
	proto "github.com/oscal-compass/compliance-to-policy/v2/api/proto/v1alpha1"
	"google.golang.org/grpc"
)

// PolicyEngine defines methods for a policy validation engine
type PolicyEngine interface {
	GetSchema() error
	UpdateConfiguration() error
	Generate() error
	GetResults() error
}

// RemediationEngine defines methods for a remediation engine
type RemediationEngine interface {
	GetSchema() error
	UpdateConfiguration() error
	Generate() error
	Remediate() error
}

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"pvp":         &PVPPlugin{},
	"remediation": &RemediationPlugin{},
}

var _ plugin.GRPCPlugin = &PVPPlugin{}
var _ plugin.GRPCPlugin = &RemediationPlugin{}

type PVPPlugin struct {
	plugin.Plugin
	Impl PolicyEngine
}

func (p *PVPPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPolicyEngineServer(s, FromPVP(p.Impl))
	return nil
}

func (p *PVPPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &pvpClient{client: proto.NewPolicyEngineClient(c)}, nil
}

type RemediationPlugin struct {
	plugin.Plugin
	Impl RemediationEngine
}

func (p *RemediationPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRemediationEngineServer(s, FromRemediation(p.Impl))
	return nil
}

func (p *RemediationPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &remediationClient{client: proto.NewRemediationEngineClient(c)}, nil
}

package plugin

import (
	"context"
	"github.com/hashicorp/go-plugin"
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/remediation"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

const (
	// PVPPluginName is used to dispense a policy plugin type
	PVPPluginName = "pvp"
	// RemediationPluginName is used to dispense a remediation plugin type
	RemediationPluginName = "remediation"
)

// Map is the map of plugins we can dispense.
var Map = map[string]plugin.Plugin{
	PVPPluginName:         &PVPPlugin{},
	RemediationPluginName: &RemediationPlugin{},
}

var _ plugin.GRPCPlugin = &PVPPlugin{}
var _ plugin.GRPCPlugin = &RemediationPlugin{}

// Below types are only used for plugins that are written in Go.

// PVPPlugin is concrete implementation of the policy Engine written in Go.
type PVPPlugin struct {
	plugin.Plugin
	Impl policy.Engine
}

func (p *PVPPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPolicyEngineServer(s, FromPVP(p.Impl))
	return nil
}

func (p *PVPPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &pvpClient{client: proto.NewPolicyEngineClient(c)}, nil
}

// RemediationPlugin is concrete implementation of the remediation Engine written in Go.
type RemediationPlugin struct {
	plugin.Plugin
	Impl remediation.Engine
}

func (p *RemediationPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRemediationEngineServer(s, FromRemediation(p.Impl))
	return nil
}

func (p *RemediationPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &remediationClient{client: proto.NewRemediationEngineClient(c)}, nil
}

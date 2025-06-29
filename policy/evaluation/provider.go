package evaluation

import (
	"context"

	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

// Provider defines methods for a policy engine C2P plugin.
type Provider interface {
	// Plan defines evaluation planning
	Plan() (*InputContext, error)
	// Report reports on evaluations
	Report(ctx context.Context, inputCtx *InputContext, output string, results []policy.PVPResult) error
}

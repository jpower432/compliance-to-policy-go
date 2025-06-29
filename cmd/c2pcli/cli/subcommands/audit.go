package subcommands

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

func NewAuditCmd(logger hclog.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "audit",
		Short: "Create audit artifacts scoped by a specific framework or standard.",
	}
	command.AddCommand(
		NewGemara2OSCAL(logger),
		NewOSCAL2Posture(logger),
	)

	return command
}

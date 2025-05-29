package eval

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

func NewCmd(logger hclog.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "eval",
		Short: "Create assessment or evaluation documents using policy-as-code artifacts.",
	}
	command.AddCommand(
		NewOSCAL2Policy(logger),
		NewResult2OSCAL(logger),
		NewResult2SCI(logger),
		NewSCI2Policy(logger),
	)

	return command
}

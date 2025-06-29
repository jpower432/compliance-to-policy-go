package subcommands

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

func NewEvalCmd(logger hclog.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "eval",
		Short: "Create assessment or evaluation documents using policy-as-code artifacts.",
	}
	command.AddCommand(
		NewOSCAL2Policy(logger),
		NewResult2Gemara(logger),
		NewGemara2Policy(logger),
	)

	return command
}

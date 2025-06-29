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
		NewResult2Compliance(logger),
		NewCompliance2Policy(logger),
	)

	return command
}

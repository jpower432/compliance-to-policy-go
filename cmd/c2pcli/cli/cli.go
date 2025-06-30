/*
Copyright 2023 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/subcommands"
	"github.com/oscal-compass/compliance-to-policy-go/v2/logging"
)

var logger hclog.Logger

func init() {
	logger = logging.GetLogger("c2pcli")
}

func New() *cobra.Command {
	var debug bool
	command := &cobra.Command{
		Use:   "c2pcli",
		Short: "C2P CLI",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				logger.SetLevel(hclog.Debug)
			}
		},
	}
	command.AddCommand(
		subcommands.NewVersionSubCommand(),
		subcommands.NewAuditCmd(logger),
		subcommands.NewEvalCmd(logger),
	)
	command.PersistentFlags().BoolVar(&debug, "debug", false, "Run with debug log level")

	return command
}

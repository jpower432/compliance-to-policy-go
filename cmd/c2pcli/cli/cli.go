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
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/subcommands"
	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/c2pcli/cli/subcommands/tools"
	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
)

var logger *zap.Logger = pkg.GetLogger("cli")

func New() *cobra.Command {

	command := &cobra.Command{
		Use:   "c2pcli",
		Short: "C2P CLI",
	}
	command.AddCommand(
		subcommands.NewVersionSubCommand(),
		tools.New(logger),
		subcommands.NewOSCAL2Policy(),
		subcommands.NewResult2OSCAL(),
	)

	return command
}

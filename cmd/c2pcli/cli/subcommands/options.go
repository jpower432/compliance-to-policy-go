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

package subcommands

import (
	"errors"

	"github.com/spf13/pflag"
)

type ResultOptions struct {
	*Options
	OutputPath string
}

func NewResultOptions(options *Options) *ResultOptions {
	return &ResultOptions{
		Options: options,
	}
}

func (o *ResultOptions) AddFlags(fs *pflag.FlagSet) {
	o.Options.AddFlags(fs)
	fs.StringVarP(&o.OutputPath, "out", "o", "./assessment-results.json", "path to output OSCAL Assessment Results")
}

func (o *ResultOptions) Validate() error {
	return o.Options.Validate()
}

type Options struct {
	ComponentDefinition string
	Name                string
	PluginsPath         string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.ComponentDefinition, "compdef", "c", ".", "path to component definition")
	fs.StringVarP(&o.Name, "name", "n", "", "short name for the chosen component implementation")
	fs.StringVarP(&o.PluginsPath, "plugin-dir", "d", "", "Path to plugin directory. Defaults to `c2p-plugins`.")
}

func (o *Options) Validate() error {
	if o.Name == "" {
		return errors.New("-n or --name is required")
	}
	return nil
}

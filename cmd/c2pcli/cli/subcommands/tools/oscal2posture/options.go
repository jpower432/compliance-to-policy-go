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

package oscal2posture

import (
	"github.com/spf13/pflag"
)

type Options struct {
	AssessmentResults   string
	Catalog             string
	ComponentDefinition string
	Out                 string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Catalog, "catalog", "c", "", "path to catalog.json")
	fs.StringVar(&o.AssessmentResults, "assessment-results", "", "path to assessment-results.json")
	fs.StringVar(&o.ComponentDefinition, "component-definition", "", "path to component-definition.json")
	fs.StringVarP(&o.Out, "out", "o", "-", "path to output file. Use '-' for stdout. Default '-'.")
}

func (o *Options) Complete() error {
	return nil
}

func (o *Options) Validate() error {
	return nil
}

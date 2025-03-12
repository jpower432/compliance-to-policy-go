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

package main

import (
	cp "github.com/otiai10/copy"

	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
	typec2pcr "github.com/oscal-compass/compliance-to-policy-go/v2/pkg/types/c2pcr"
)

type Oscal2Policy struct {
}

func NewOscal2Policy(policiesDir string, tempDir pkg.TempDirectory) *Oscal2Policy {
	return &Oscal2Policy{
		policiesDir: policiesDir,
		tempDir:     tempDir,
		logger:      pkg.GetLogger("kyverno/composer"),
	}
}

func (c *Oscal2Policy) Generate(c2pParsed typec2pcr.C2PCRParsed) error {

	return nil
}

func (c *Oscal2Policy) CopyAllTo(destDir string) error {
	if _, err := pkg.MakeDir(destDir); err != nil {
		return err
	}
	if err := cp.Copy(c.tempDir.GetTempDir(), destDir); err != nil {
		return err
	}
	return nil
}

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

package parser

import (
	"os"
	"path/filepath"
)

func WriteToCSVs(c *Collector, outputDir string) (string, string) {
	policyCsvPath := outputDir + "/policies.csv"
	policyCsvPath = filepath.Clean(policyCsvPath)
	of, err := os.Create(policyCsvPath)
	if err != nil {
		panic(err)
	}
	c.GetTable().ToCsv(of)

	resourcesCsvPath := outputDir + "/resources.csv"
	resourcesCsvPath = filepath.Clean(resourcesCsvPath)

	of, err = os.Create(resourcesCsvPath)
	if err != nil {
		panic(err)
	}
	c.GetResourceTable().ToCsv(of)
	return policyCsvPath, resourcesCsvPath
}

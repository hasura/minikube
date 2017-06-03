/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"
	"github.com/hasura/minikube/cmd/minikube/cmd"
)

func main() {
	os.MkdirAll("./out/docs", os.FileMode(0755))
	cmd.RootCmd.DisableAutoGenTag = true
	doc.GenMarkdownTree(cmd.RootCmd, "./out/docs")

	f, err := os.Create("./out/docs/bash-completion")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	err = cmd.GenerateBashCompletion(f, cmd.RootCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

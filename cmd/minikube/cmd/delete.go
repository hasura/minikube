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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cmdUtil "gitlab.com/hasura/hasuractl-go/pkg/minikube/cmd/util"
	"gitlab.com/hasura/hasuractl-go/pkg/minikube/pkg/minikube/cluster"
	"gitlab.com/hasura/hasuractl-go/pkg/minikube/pkg/minikube/machine"
)

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a local kubernetes cluster",
	Long: `Deletes a local kubernetes cluster. This command deletes the VM, and removes all
associated files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting local Kubernetes cluster...")
		api, err := machine.NewAPIClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting client: %s\n", err)
			os.Exit(1)
		}
		defer api.Close()

		if err = cluster.DeleteHost(api); err != nil {
			fmt.Println("Errors occurred deleting machine: ", err)
			os.Exit(1)
		}
		fmt.Println("Machine deleted.")

		if err := cmdUtil.KillMountProcess(); err != nil {
			fmt.Println("Errors occurred deleting mount process: ", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
}

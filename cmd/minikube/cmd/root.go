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
	goflag "flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/docker/machine/libmachine/log"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	configCmd "gitlab.com/hasura/hasuractl-go/pkg/minikube/cmd/minikube/cmd/config"
	"gitlab.com/hasura/hasuractl-go/pkg/minikube/cmd/util"
	"gitlab.com/hasura/hasuractl-go/pkg/minikube/pkg/minikube/config"
	"gitlab.com/hasura/hasuractl-go/pkg/minikube/pkg/minikube/constants"
	"gitlab.com/hasura/hasuractl-go/pkg/minikube/pkg/minikube/notify"
)

var dirs = [...]string{
	constants.GetMinipath(),
	constants.MakeMiniPath("certs"),
	constants.MakeMiniPath("machines"),
	constants.MakeMiniPath("cache"),
	constants.MakeMiniPath("cache", "iso"),
	constants.MakeMiniPath("cache", "localkube"),
	constants.MakeMiniPath("config"),
	constants.MakeMiniPath("addons"),
	constants.MakeMiniPath("logs"),
}

const (
	showLibmachineLogs = "show-libmachine-logs"
	useVendoredDriver  = "use-vendored-driver"
)

var (
	enableUpdateNotification = true
)

var viperWhiteList = []string{
	"v",
	"alsologtostderr",
	"log_dir",
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "minikube",
	Short: "Minikube is a tool for managing local Kubernetes clusters.",
	Long:  `Minikube is a CLI tool that provisions and manages single-node Kubernetes clusters optimized for development workflows.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		for _, path := range dirs {
			if err := os.MkdirAll(path, 0777); err != nil {
				glog.Exitf("Error creating minikube directory: %s", err)
			}
		}

		// Log level 3 or greater enables libmachine logs
		if !glog.V(3) {
			log.SetOutWriter(ioutil.Discard)
			log.SetErrWriter(ioutil.Discard)
		}

		// Log level 7 or greater enables debug level logs
		if glog.V(7) {
			log.SetDebug(true)
		}

		if viper.GetBool(showLibmachineLogs) {
			fmt.Println(`
--show-libmachine-logs is deprecated.
Please use --v=3 to show libmachine logs, and --v=7 for debug level libmachine logs
`)
		}

		logDir := pflag.Lookup("log_dir")
		if !logDir.Changed {
			logDir.Value.Set(constants.MakeMiniPath("logs"))
		}

		if enableUpdateNotification {
			notify.MaybePrintUpdateTextFromGithub(os.Stderr)
		}
		util.MaybePrintKubectlDownloadMsg(runtime.GOOS, os.Stderr)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_ = RootCmd.Execute()
}

// Handle config values for flags used in external packages (e.g. glog)
// by setting them directly, using values from viper when not passed in as args
func setFlagsUsingViper() {
	for _, config := range viperWhiteList {
		var a = pflag.Lookup(config)
		viper.SetDefault(a.Name, a.DefValue)
		// If the flag is set, override viper value
		if a.Changed {
			viper.Set(a.Name, a.Value.String())
		}
		// Viper will give precedence first to calls to the Set command,
		// then to values from the config.yml
		a.Value.Set(viper.GetString(a.Name))
		a.Changed = true
	}
}

func init() {
	RootCmd.PersistentFlags().Bool(showLibmachineLogs, false, "Deprecated: To enable libmachine logs, set --v=3 or higher")
	RootCmd.PersistentFlags().Bool(useVendoredDriver, false, "Use the vendored in drivers instead of RPC")
	RootCmd.PersistentFlags().StringP(config.MachineProfile, "p", constants.DefaultMachineName, `The name of the minikube VM being used.  
	This can be modified to allow for multiple minikube instances to be run independently`)
	RootCmd.AddCommand(configCmd.ConfigCmd)
	RootCmd.AddCommand(configCmd.AddonsCmd)
	RootCmd.AddCommand(configCmd.ProfileCmd)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	viper.BindPFlags(RootCmd.PersistentFlags())

	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	configPath := constants.ConfigFile
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		glog.Warningf("Error reading config file at %s: %s", configPath, err)
	}
	setupViper()
}

func setupViper() {
	viper.SetEnvPrefix(constants.MinikubeEnvPrefix)
	// Replaces '-' in flags with '_' in env variables
	// e.g. iso-url => $ENVPREFIX_ISO_URL
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetDefault(config.WantUpdateNotification, true)
	viper.SetDefault(config.ReminderWaitPeriodInHours, 24)
	viper.SetDefault(config.WantReportError, false)
	viper.SetDefault(config.WantReportErrorPrompt, true)
	viper.SetDefault(config.WantKubectlDownloadMsg, true)
	setFlagsUsingViper()
}

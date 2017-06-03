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

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/machine/libmachine/drivers"

	"github.com/hasura/minikube/pkg/minikube/assets"
	"github.com/hasura/minikube/pkg/minikube/sshutil"
	"github.com/hasura/minikube/pkg/minikube/tests"
)

func TestDisableUnknownAddon(t *testing.T) {
	if err := Set("InvalidAddon", "false"); err == nil {
		t.Fatalf("Disable did not return error for unknown addon")
	}
}

func TestDisableValidAddonNoVM(t *testing.T) {
	tempDir := tests.MakeTempDir()
	defer os.RemoveAll(tempDir)

	addonName := "dashboard"
	expected := "false"
	if err := Set(addonName, expected); err != nil {
		//check that setting was changed in temp dir
		conf, _ := ioutil.ReadFile(filepath.Join(tempDir, "config", "config.json"))
		strConf := string(conf)
		if !strings.Contains(strConf, addonName) {
			t.Fatalf("Expected %s, in %s after running disable", addonName, strConf)
		}
		if !strings.Contains(strConf, expected) {
			t.Fatalf("Expected %s, in %s after running disable", expected, strConf)
		}
		//check that nothing was transferred?
	} else {
		if err != nil {
			t.Fatalf("Disable returned error for known addon")
		}
	}
}

func TestDeleteAddonViaDriver(t *testing.T) {
	s, _ := tests.NewSSHServer()
	port, err := s.Start()
	if err != nil {
		t.Fatalf("Error starting ssh server: %s", err)
	}

	d := &tests.MockDriver{
		Port: port,
		BaseDriver: drivers.BaseDriver{
			IPAddress:  "127.0.0.1",
			SSHKeyPath: "",
		},
	}

	dashboard := assets.Addons["dashboard"]
	if err := deleteAddonViaDriver(dashboard, d); err != nil {
		t.Fatalf("Unexpected error %s deleting addon", err)
	}
	// check command(s) were run
	for _, addon := range dashboard.Assets {
		expected, _ := ioutil.ReadFile(addon.GetAssetName())
		if _, ok := s.Commands[sshutil.GetDeleteFileCommand(addon)]; !ok {
			t.Fatalf("Error: Expected delete addon ssh command to be run: %s.", expected)
		}
	}
}

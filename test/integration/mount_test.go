// +build integration

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

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"k8s.io/kubernetes/pkg/api"
	commonutil "github.com/hasura/minikube/pkg/util"
	"github.com/hasura/minikube/test/integration/util"
)

func testMounting(t *testing.T) {
	t.Parallel()
	minikubeRunner := util.MinikubeRunner{
		Args:       *args,
		BinaryPath: *binaryPath,
		T:          t}

	tempDir, err := ioutil.TempDir("", "mounttest")
	if err != nil {
		t.Fatalf("Unexpected error while creating tempDir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	mountCmd := fmt.Sprintf("mount %s", tempDir)
	cmd := minikubeRunner.RunDaemon(mountCmd)
	defer cmd.Process.Kill()

	kubectlRunner := util.NewKubectlRunner(t)
	podName := "busybox"
	podPath, _ := filepath.Abs("testdata/busybox-mount-test.yaml")

	// Write file in mounted dir from host
	expected := "test\n"
	files := []string{"fromhost", "fromhostremove"}
	for _, file := range files {
		path := filepath.Join(tempDir, file)
		err = ioutil.WriteFile(path, []byte(expected), 0644)
		if err != nil {
			t.Fatalf("Unexpected error while writing file %s: %s.", path, err)
		}
	}
	mountTest := func() error {
		if _, err := kubectlRunner.RunCommand([]string{"create", "-f", podPath}); err != nil {
			return err
		}
		defer kubectlRunner.RunCommand([]string{"delete", "-f", podPath})

		p := &api.Pod{}
		for p.Status.Phase != "Running" {
			p = kubectlRunner.GetPod(podName, "default")
		}

		path := filepath.Join(tempDir, "frompod")
		out, err := ioutil.ReadFile(path)
		if err != nil {
			return &commonutil.RetriableError{Err: err}
		}
		// test that file written from pod can be read from host echo test > /mount-9p/frompod; in pod
		if string(out) != expected {
			t.Fatalf("Expected file %s to contain text %s, was %s.", path, expected, out)
		}

		// test that file written from host was read in by the pod via cat /mount-9p/fromhost;
		if out, err = kubectlRunner.RunCommand([]string{"logs", podName}); err != nil {
			return &commonutil.RetriableError{Err: err}
		}
		if string(out) != expected {
			t.Fatalf("Expected file %s to contain text %s, was %s.", path, expected, out)
		}

		// test that fromhostremove was deleted by the pod from the mount via rm /mount-9p/fromhostremove
		path = filepath.Join(tempDir, "fromhostremove")
		if _, err := os.Stat(path); err == nil {
			t.Fatalf("Expected file %s to be removed", path, expected, out)
		}

		// test that frompodremove can be deleted on the host
		path = filepath.Join(tempDir, "frompodremove")
		if err := os.Remove(path); err != nil {
			t.Fatalf("Unexpected error removing file %s: %s", path, err)
		}

		return nil
	}
	if err := commonutil.RetryAfter(40, mountTest, 5*time.Second); err != nil {
		t.Fatal("mountTest failed with error:", err)
	}

}

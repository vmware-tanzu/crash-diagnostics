// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestKubeExecScript(t *testing.T) {
	workdir := testSupport.TmpDirRoot()
	k8sconfig := testSupport.KindKubeConfigFile()
	clusterName := testSupport.KindClusterContextName()
	err := testSupport.StartNginxPod()
	if err != nil {
		t.Error("Unexpected error while starting nginx pod", err)
		return
	}

	execute := func(t *testing.T, script string) *starlarkstruct.Struct {
		executor := New()
		if err := executor.Exec("test.kube.exec", strings.NewReader(script)); err != nil {
			t.Fatalf("failed to exec: %s", err)
		}
		if !executor.result.Has("kube_exec_output") {
			t.Fatalf("script result must be assigned to a value")
		}

		data, ok := executor.result["kube_exec_output"].(*starlarkstruct.Struct)
		if !ok {
			t.Fatal("script result is not a struct")
		}
		return data
	}

	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "exec into pod and run long-running operation",
			script: fmt.Sprintf(`
crashd_config(workdir="%s")
set_defaults(kube_config(path="%s", cluster_context="%s"))
kube_exec_output=kube_exec(pod="nginx", timeout_in_seconds=3,cmd=["sh", "-c" ,"while true; do echo 'Running'; sleep 1; done"])
`, workdir, k8sconfig, clusterName),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				errVal, err := data.Attr("error")
				if err != nil {
					t.Fatal(err)
				}

				resultErr := errVal.(starlark.String).GoString()
				if resultErr == "" || !strings.HasPrefix(resultErr, "command execution timed out.") {
					t.Fatalf("Unexpected error result: %s", resultErr)
				}

				ouputFilePath, err := data.Attr("file")
				if err != nil {
					t.Fatal(err)
				}

				file, err := os.Open(trimQuotes(ouputFilePath.String()))
				if err != nil {
					t.Fatalf("result file does not exist: %s", err)
				}
				defer file.Close()

				var actual int
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if line == "Running" {
						actual++
					}
				}
				expected := 3
				if expected != actual {
					t.Fatalf("Unexpected file content. expected line numbers: %d but was %d", expected, actual)
				}

			},
		},
		{
			name: "exec into pod and run short-lived command.Output to specified file",
			script: fmt.Sprintf(`
crashd_config(workdir="%s")
set_defaults(kube_config(path="%s", cluster_context="%s"))
kube_exec_output=kube_exec(pod="nginx", output_file="nginx.version",container="nginx", cmd=["nginx", "-v"])
`, workdir, k8sconfig, clusterName),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				errVal, err := data.Attr("error")
				if err != nil {
					t.Fatal(err)
				}

				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("expected ouput error to be empty but was %s", resultErr)
				}

				ouputFilePath, err := data.Attr("file")
				if err != nil {
					t.Fatal(err)
				}

				fileContents, err := os.ReadFile(trimQuotes(ouputFilePath.String()))
				if err != nil {
					t.Fatalf("Error reading output file: %v", err)
				}
				strings.Contains(string(fileContents), "nginx version:")

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script)
		})
	}
}

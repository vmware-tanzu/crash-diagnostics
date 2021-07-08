// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script_tests

import (
	"os"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/functions/archive"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

func TestArchiveScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "archive defaults",
			script: `
result = archive(output_file="/tmp/archive.tar.gz", source_paths=["/tmp/crashd"])
`,
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}

				expected := "/tmp/archive.tar.gz"
				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("archive() should be assigned to a variable for test")
				}
				var result archive.Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}

				defer func() {
					os.RemoveAll(expected)
					os.RemoveAll("/tmp/crashd")
				}()

				if result.Archive.OutputFile != expected {
					t.Errorf("unexpected result: %s", result.Archive.OutputFile)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script)
		})
	}
}

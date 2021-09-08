//// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
//// SPDX-License-Identifier: Apache-2.0
//
package archive

import (
	"os"
	"strings"
	"testing"

	"go.starlark.net/starlark"

	crashlark "github.com/vmware-tanzu/crash-diagnostics/starlark"
	crashtest "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestArchiveFunc(t *testing.T) {
	tests := []struct {
		name string
		args func(t *testing.T) []starlark.Tuple
		eval func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name: "arhive single file",
			args: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					{starlark.String("output_file"), starlark.String("/tmp/out.tar.gz")},
					{starlark.String("source_paths"), starlark.NewList([]starlark.Value{starlark.String("/tmp/crashd")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := archiveFunc(crashtest.NewStarlarkThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				expected := "/tmp/out.tar.gz"
				defer func() {
					os.RemoveAll(expected)
					os.RemoveAll("/tmp/crashd")
				}()

				result := ""
				if r, ok := val.(starlark.String); ok {
					result = string(r)
				}
				if result != expected {
					t.Errorf("unexpected result: %s", result)
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.args(t))
		})
	}
}

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
				exe := crashlark.New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				expected := "/tmp/archive.tar.gz"
				var result string
				resultVal := exe.Result()["result"]
				if resultVal == nil {
					t.Fatal("archive() should be assigned to a variable for test")
				}
				res, ok := resultVal.(starlark.String)
				if !ok {
					t.Fatal("archive() should return a string")
				}
				result = string(res)
				defer func() {
					os.RemoveAll(result)
					os.RemoveAll("/tmp/crashd")
				}()

				if result != expected {
					t.Errorf("unexpected result: %s", result)
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

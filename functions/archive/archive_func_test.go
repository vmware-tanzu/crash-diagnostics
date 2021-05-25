//// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
//// SPDX-License-Identifier: Apache-2.0
//
package archive

import (
	"errors"
	"os"
	"strings"
	"testing"

	"go.starlark.net/starlark"

	crashlark "github.com/vmware-tanzu/crash-diagnostics/starlark"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

func TestArchiveFunc(t *testing.T) {
	tests := []struct {
		name string
		args func(t *testing.T) []starlark.Tuple
		eval func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name: "archive single file",
			args: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					{starlark.String("output_file"), starlark.String("/tmp/out.tar.gz")},
					{starlark.String("source_paths"), starlark.NewList([]starlark.Value{starlark.String("/tmp/crashd")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				if err := createSourceFile(t, "/tmp/crashd", "test.txt", "Hello"); err != nil {
					t.Fatal(err)
				}

				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				var arc Archive
				if err := typekit.Starlark(val).Go(&arc); err != nil {
					t.Fatal(err)
				}

				if arc.OutputFile != "/tmp/out.tar.gz" {
					t.Errorf("unexpected output file: %s", arc.OutputFile)
				}
				if _, err := os.Stat(arc.OutputFile); err != nil {
					t.Fatal(err)
				}

				if len(arc.SourcePaths) != 1 {
					t.Errorf("unexpected source paths len: %d", len(arc.SourcePaths))
				}

				if err := os.RemoveAll(arc.OutputFile); err != nil {
					t.Log(err)
				}
				for _, p := range arc.SourcePaths {
					if err := os.RemoveAll(p); err != nil {
						t.Log(err)
					}
				}
			},
		},
		{
			name: "archive multiple files",
			args: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					{starlark.String("output_file"), starlark.String("/tmp/out.tar.gz")},
					{starlark.String("source_paths"), starlark.NewList(
						[]starlark.Value{starlark.String("/tmp/crashd")},
					)},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				if err := createSourceFile(t, "/tmp/crashd", "1.txt", "Hello"); err != nil {
					t.Fatal(err)
				}
				if err := createSourceFile(t, "/tmp/crashd", "2.txt", "Hello"); err != nil {
					t.Fatal(err)
				}
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				var arc Archive
				if err := typekit.Starlark(val).Go(&arc); err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(arc.OutputFile); err != nil {
					t.Fatal(err)
				}

				if err := os.RemoveAll(arc.OutputFile); err != nil {
					t.Log(err)
				}
				for _, p := range arc.SourcePaths {
					if err := os.RemoveAll(p); err != nil {
						t.Log(err)
					}
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

func createSourceFile(t *testing.T, rootPath, fileName, content string) error {
	t.Logf("creating archive source: %s/%s", rootPath, fileName)
	if err := os.MkdirAll(rootPath, 0754); !errors.Is(err, os.ErrExist) {
		return err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

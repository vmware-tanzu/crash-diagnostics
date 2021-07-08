// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package archive

import (
	"errors"
	"os"
	"testing"

	"go.starlark.net/starlark"

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
				defer os.RemoveAll("/tmp/crashd")

				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				var result Result
				if err := typekit.Starlark(val).Go(&result); err != nil {
					t.Fatal(err)
				}

				if result.Archive.OutputFile != "/tmp/out.tar.gz" {
					t.Errorf("unexpected output file: %s", result.Archive.OutputFile)
				}
				if _, err := os.Stat(result.Archive.OutputFile); err != nil {
					t.Fatal(err)
				}

				if err := os.RemoveAll(result.Archive.OutputFile); err != nil {
					t.Log(err)
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
				defer os.RemoveAll("/tmp/crashd")

				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				var result Result
				if err := typekit.Starlark(val).Go(&result); err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(result.Archive.OutputFile); err != nil {
					t.Fatal(err)
				}

				if err := os.RemoveAll(result.Archive.OutputFile); err != nil {
					t.Log(err)
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

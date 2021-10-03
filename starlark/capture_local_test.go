// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestCaptureLocalFunc(t *testing.T) {
	tests := []struct {
		name string
		args func(t *testing.T) []starlark.Tuple
		eval func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name: "capture with defaults",
			args: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{{starlark.String("cmd"), starlark.String("echo 'Hello World!'")}}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := captureLocalFunc(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				expected := filepath.Join(defaults.workdir, fmt.Sprintf("%s.txt", sanitizeStr("echo 'Hello World!'")))
				result := ""
				if r, ok := val.(starlark.String); ok {
					result = string(r)
				}
				defer func() {
					os.RemoveAll(result)
					os.RemoveAll(defaults.workdir)
				}()

				if result != expected {
					t.Errorf("unexpected result: %s", result)
				}

				file, err := os.Open(result)
				if err != nil {
					t.Fatal(err)
				}
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, file); err != nil {
					t.Fatal(err)
				}
				expected = strings.TrimSpace(buf.String())
				if expected != "Hello World!" {
					t.Errorf("unexpected content captured: %s", expected)
				}
				if err := file.Close(); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "capture with args",
			args: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					{starlark.String("workdir"), starlark.String("/tmp/capturecrashd")},
					{starlark.String("file_name"), starlark.String("echo.txt")},
					{starlark.String("desc"), starlark.String("echo command")},
					{starlark.String("append"), starlark.True},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				expected := filepath.Join("/tmp/capturecrashd", "echo.txt")
				err := os.MkdirAll("/tmp/capturecrashd", 0777)
				if err != nil {
					t.Fatal(err)
				}
				f, err := os.OpenFile(expected, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := f.Write([]byte("Hello World!\n")); err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
				val, err := captureLocalFunc(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				result := ""
				if r, ok := val.(starlark.String); ok {
					result = string(r)
				}
				defer func() {
					os.RemoveAll(result)
					os.RemoveAll(defaults.workdir)
				}()

				if result != expected {
					t.Errorf("unexpected result: %s", result)
				}

				file, err := os.Open(result)
				if err != nil {
					t.Fatal(err)
				}
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, file); err != nil {
					t.Fatal(err)
				}
				expected = strings.TrimSpace(buf.String())
				if expected != "Hello World!\necho command\nHello World!" {
					t.Errorf("unexpected content captured: %s", expected)
				}
				if err := file.Close(); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "capture with error",
			args: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{{starlark.String("cmd"), starlark.String("nacho 'Hello World!'")}}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := captureLocalFunc(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				result := ""
				if r, ok := val.(starlark.String); ok {
					result = string(r)
				}
				defer func() {
					os.RemoveAll(result)
					os.RemoveAll(defaults.workdir)
				}()

				file, err := os.Open(result)
				if err != nil {
					t.Fatal(err)
				}
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, file); err != nil {
					t.Fatal(err)
				}
				expected := strings.TrimSpace(buf.String())
				if !strings.Contains(expected, "not found") {
					t.Errorf("unexpected content captured: %s", expected)
				}
				if err := file.Close(); err != nil {
					t.Error(err)
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

func TestCaptureLocalScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "capture local defaults",
			script: `
result = capture_local("echo 'Hello World!'")
`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				expected := filepath.Join(defaults.workdir, fmt.Sprintf("%s.txt", sanitizeStr("echo 'Hello World!'")))
				var result string
				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("capture_local() should be assigned to a variable for test")
				}
				res, ok := resultVal.(starlark.String)
				if !ok {
					t.Fatal("capture_local() should return a string")
				}
				result = string(res)
				defer func() {
					os.RemoveAll(result)
					os.RemoveAll(defaults.workdir)
				}()

				if result != expected {
					t.Errorf("unexpected result: %s", result)
				}

				file, err := os.Open(result)
				if err != nil {
					t.Fatal(err)
				}
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, file); err != nil {
					t.Fatal(err)
				}
				expected = strings.TrimSpace(buf.String())
				if expected != "Hello World!" {
					t.Errorf("unexpected content captured: %s", expected)
				}
				if err := file.Close(); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "capture local with args",
			script: `
result = capture_local(cmd="echo 'Hello World!'", workdir="/tmp/capturecrash", file_name="echo_out.txt", desc="output command", append=True)
`,
			eval: func(t *testing.T, script string) {
				expected := filepath.Join("/tmp/capturecrash", "echo_out.txt")
				err := os.MkdirAll("/tmp/capturecrash", 0777)
				if err != nil {
					t.Fatal(err)
				}
				f, err := os.OpenFile(expected, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := f.Write([]byte("Hello World!\n")); err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				var result string
				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("capture_local() should be assigned to a variable for test")
				}
				res, ok := resultVal.(starlark.String)
				if !ok {
					t.Fatal("capture_local() should return a string")
				}
				result = string(res)
				defer func() {
					os.RemoveAll(result)
					os.RemoveAll(defaults.workdir)
				}()

				if result != expected {
					t.Errorf("unexpected result: %s", result)
				}

				file, err := os.Open(result)
				if err != nil {
					t.Fatal(err)
				}
				buf := new(bytes.Buffer)
				if _, err := io.Copy(buf, file); err != nil {
					t.Fatal(err)
				}
				expected = strings.TrimSpace(buf.String())
				if expected != "Hello World!\noutput command\nHello World!" {
					t.Errorf("unexpected content captured: %s", expected)
				}
				if err := file.Close(); err != nil {
					t.Error(err)
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

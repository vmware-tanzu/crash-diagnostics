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

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func testCaptureFuncForHostResources(t *testing.T, port string) {
	tests := []struct {
		name   string
		args   func(t *testing.T) starlark.Tuple
		kwargs func(t *testing.T) []starlark.Tuple
		eval   func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple)
	}{
		{
			name: "default args single machine",
			args: func(t *testing.T) starlark.Tuple { return starlark.Tuple{starlark.String("echo 'Hello World!'")} },
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(defaults.pkPath, port)
				resources := starlark.NewList([]starlark.Value{makeTestSSHHostResource("127.0.0.1", sshCfg)})
				return []starlark.Tuple{[]starlark.Value{starlark.String("resources"), resources}}
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				val, err := captureFunc(newTestThreadLocal(t), nil, args, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				result := ""
				if strct, ok := val.(*starlarkstruct.Struct); ok {
					if val, err := strct.Attr("result"); err == nil {
						if r, ok := val.(starlark.String); ok {
							result = string(r)
						}
					}
				}

				expected := filepath.Join(defaults.workdir, sanitizeStr("127.0.0.1"), fmt.Sprintf("%s.txt", sanitizeStr("echo 'Hello World!'")))
				if result != expected {
					t.Errorf("unexpected file name captured: %s", result)
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
				defer os.RemoveAll(result)
			},
		},

		{
			name: "kwargs single machine",
			args: func(t *testing.T) starlark.Tuple { return nil },
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(defaults.pkPath, port)
				resources := starlark.NewList([]starlark.Value{makeTestSSHHostResource("127.0.0.1", sshCfg)})
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					[]starlark.Value{starlark.String("resources"), resources},
					[]starlark.Value{starlark.String("file_name"), starlark.String("echo_out.txt")},
					[]starlark.Value{starlark.String("desc"), starlark.String("echo command")},
				}
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				val, err := captureFunc(newTestThreadLocal(t), nil, args, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				result := ""
				if strct, ok := val.(*starlarkstruct.Struct); ok {
					if val, err := strct.Attr("result"); err == nil {
						if r, ok := val.(starlark.String); ok {
							result = string(r)
						}
					}
				}
				expected := filepath.Join(defaults.workdir, sanitizeStr("127.0.0.1"), "echo_out.txt")
				if result != expected {
					t.Errorf("unexpected file name captured: %s", result)
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
				if expected != "echo command\nHello World!" {
					t.Errorf("unexpected content captured: %s", expected)
				}
				if err := file.Close(); err != nil {
					t.Error(err)
				}
				defer os.RemoveAll(result)
			},
		},

		{
			name: "multiple machines",
			args: func(t *testing.T) starlark.Tuple { return nil },
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(defaults.pkPath, port)
				resources := starlark.NewList([]starlark.Value{
					makeTestSSHHostResource("localhost", sshCfg),
					makeTestSSHHostResource("127.0.0.1", sshCfg),
				})
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					[]starlark.Value{starlark.String("resources"), resources},
				}
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				val, err := captureFunc(newTestThreadLocal(t), nil, args, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				resultList, ok := val.(*starlark.List)
				if !ok {
					t.Fatalf("expecting type *starlark.List, got %T", val)
				}

				for i := 0; i < resultList.Len(); i++ {
					result := ""
					if strct, ok := resultList.Index(i).(*starlarkstruct.Struct); ok {
						if val, err := strct.Attr("result"); err == nil {
							if r, ok := val.(starlark.String); ok {
								result = string(r)
							}
						}
					}
					if _, err := os.Stat(result); err != nil {
						t.Fatalf("captured command file not found: %s", err)
					}
					os.RemoveAll(result)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.args(t), test.kwargs(t))
		})
	}
}

func testCaptureFuncScriptForHostResources(t *testing.T, port string) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "default cmd multiple machines",
			script: fmt.Sprintf(`
set_defaults(resources(provider = host_list_provider(hosts=["127.0.0.1","localhost"], ssh_config = ssh_config(username=os.username, port="%s"))))
result = capture("echo 'Hello World!'")`, port),
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("capture() should be assigned to a variable")
				}
				resultList, ok := resultVal.(*starlark.List)
				if !ok {
					t.Fatal("capture() with multiple resources should return a list")
				}

				for i := 0; i < resultList.Len(); i++ {
					resultStruct, ok := resultList.Index(i).(*starlarkstruct.Struct)
					if !ok {
						t.Fatalf("capture(): expecting a starlark struct, got %T", resultList.Index(i))
					}
					val, err := resultStruct.Attr("result")
					if err != nil {
						t.Fatal(err)
					}
					result := string(val.(starlark.String))
					if _, err := os.Stat(result); err != nil {
						t.Fatalf("captured command file not found: %s", err)
					}
					os.RemoveAll(result)
				}
			},
		},

		{
			name: "resource loop",
			script: fmt.Sprintf(`
# execute cmd on each host
def exec(hosts):
	result = []
	for host in hosts:
		result.append(capture(cmd="echo 'Hello World!'", resources=[host], file_name="echo.txt", desc="echo command:"))
	return result
		
# configuration
set_defaults(ssh_config(username=os.username, port="%s"))
hosts = resources(provider=host_list_provider(hosts=["127.0.0.1","localhost"]))
result = exec(hosts)`, port),
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("capture() should be assigned to a variable")
				}
				resultList, ok := resultVal.(*starlark.List)
				if !ok {
					t.Fatal("capture() with multiple resources should return a list")
				}

				for i := 0; i < resultList.Len(); i++ {
					resultStruct, ok := resultList.Index(i).(*starlarkstruct.Struct)
					if !ok {
						t.Fatalf("run(): expecting a starlark struct, got %T", resultList.Index(i))
					}
					val, err := resultStruct.Attr("result")
					if err != nil {
						t.Fatal(err)
					}
					result := string(val.(starlark.String))
					if _, err := os.Stat(result); err != nil {
						t.Fatalf("captured command file not found: %s", err)
					}
					//os.RemoveAll(result)
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

func TestCaptureFuncSSHAll(t *testing.T) {
	port := testcrashd.NextPortValue()
	sshSvr := testcrashd.NewSSHServer(testcrashd.NextResourceName(), port)

	logrus.Debug("Attempting to start SSH server")
	if err := sshSvr.Start(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	tests := []struct {
		name string
		test func(t *testing.T, port string)
	}{
		{name: "capture func for host resources", test: testCaptureFuncForHostResources},
		{name: "capture script for host resources", test: testCaptureFuncScriptForHostResources},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.test(t, port)
			defer os.RemoveAll(defaults.workdir)
		})
	}

	logrus.Debug("Stopping SSH server...")
	if err := sshSvr.Stop(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func testRunFuncHostResources(t *testing.T, port string) {
	tests := []struct {
		name   string
		args   func(t *testing.T) starlark.Tuple
		kwargs func(t *testing.T) []starlark.Tuple
		eval   func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple)
	}{
		{
			name: "default arg single machine",
			args: func(t *testing.T) starlark.Tuple { return starlark.Tuple{starlark.String("echo 'Hello World!'")} },
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(testcrashd.GetSSHPrivateKey(), port)
				resources := starlark.NewList([]starlark.Value{makeTestSSHHostResource("127.0.0.1", sshCfg)})
				return []starlark.Tuple{[]starlark.Value{starlark.String("resources"), resources}}
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				val, err := runFunc(newTestThreadLocal(t), nil, args, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				expected := "Hello World!"
				result := ""
				if strct, ok := val.(*starlarkstruct.Struct); ok {
					if val, err := strct.Attr("result"); err == nil {
						if r, ok := val.(starlark.String); ok {
							result = string(r)
						}
					}
				}
				if expected != result {
					t.Fatalf("runFunc returned unexpected value: %s", string(val.(starlark.String)))
				}
			},
		},

		{
			name: "kwargs single machine",
			args: func(t *testing.T) starlark.Tuple { return nil },
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(testcrashd.GetSSHPrivateKey(), port)
				resources := starlark.NewList([]starlark.Value{makeTestSSHHostResource("127.0.0.1", sshCfg)})
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					[]starlark.Value{starlark.String("resources"), resources},
				}
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				val, err := runFunc(newTestThreadLocal(t), nil, args, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				expected := "Hello World!"
				result := ""
				if strct, ok := val.(*starlarkstruct.Struct); ok {
					if val, err := strct.Attr("result"); err == nil {
						if r, ok := val.(starlark.String); ok {
							result = string(r)
						}
					}
				}
				if expected != result {
					t.Fatalf("runFunc returned unexpected value: %s", string(val.(starlark.String)))
				}
			},
		},

		{
			name: "multiple machines",
			args: func(t *testing.T) starlark.Tuple { return nil },
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(testcrashd.GetSSHPrivateKey(), port)
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
				val, err := runFunc(newTestThreadLocal(t), nil, args, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				resultList, ok := val.(*starlark.List)
				if !ok {
					t.Fatalf("expecting type *starlark.List, got %T", val)
				}

				for i := 0; i < resultList.Len(); i++ {
					expected := "Hello World!"
					result := ""
					if strct, ok := resultList.Index(i).(*starlarkstruct.Struct); ok {
						if val, err := strct.Attr("result"); err == nil {
							if r, ok := val.(starlark.String); ok {
								result = string(r)
							}
						}
					}
					if expected != result {
						t.Fatalf("runFunc returned unexpected value: %s", string(val.(starlark.String)))
					}
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

func testRunFuncScriptHostResources(t *testing.T, port string) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "default cmd multiple machines",
			script: fmt.Sprintf(`
set_as_default(ssh_config = ssh_config(username="%s", port="%s", private_key_path="%s"))
set_as_default(resources = resources(hosts=["127.0.0.1","localhost"]))
result = run("echo 'Hello World!'")`, testcrashd.GetSSHUsername(), port, testcrashd.GetSSHPrivateKey()),
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("run() should be assigned to a variable")
				}
				resultList, ok := resultVal.(*starlark.List)
				if !ok {
					t.Fatal("run() with multiple resources should return a list")
				}
				expected := "Hello World!"
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
					if expected != result {
						t.Errorf("run(): expecting %s, got %s", expected, result)
					}
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
		result.append(run(cmd="echo 'Hello World!'", resources=[host]))
	return result

# configuration
hosts = resources(provider=host_list_provider(hosts=["127.0.0.1","localhost"], ssh_config = ssh_config(username="%s", port="%s", private_key_path="%s")))
result = exec(hosts)`, testcrashd.GetSSHUsername(), port, testcrashd.GetSSHPrivateKey()),
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("run() should be assigned to a variable")
				}
				resultList, ok := resultVal.(*starlark.List)
				if !ok {
					t.Fatal("run() with multiple resources should return a list")
				}
				expected := "Hello World!"
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
					if expected != result {
						t.Errorf("run(): expecting %s, got %s", expected, result)
					}
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

func TestRunFuncSSHAll(t *testing.T) {
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
		{name: "testRunFuncWithHostResources", test: testRunFuncHostResources},
		{name: "testRunFuncScriptWithHostResources", test: testRunFuncScriptHostResources},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { test.test(t, port) })
	}

	logrus.Debug("Stopping SSH server...")
	if err := sshSvr.Stop(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

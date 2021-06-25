// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers/hostlist"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestRunFunc(t *testing.T) {
	tests := []struct {
		name   string
		kwargs func(*testing.T) []starlark.Tuple
		eval   func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name: "no args",
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				t.Logf("expected err: %s", err)
				if err == nil {
					t.Fatal("expected error, but got nil")
				}

			},
		},
		{
			name: "missing resources",
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshConf := sshconf.DefaultConfig()
				sshArg, err := functions.Result(sshconf.Name, sshConf)
				if err != nil {
					t.Fatal(err)
				}
				return []starlark.Tuple{
					{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				t.Logf("expected err: %s", err)
				if err == nil {
					t.Fatal("expecting error, but got nil")
				}
			},
		},
		{
			name: "missing sshconf",
			kwargs: func(t *testing.T) []starlark.Tuple {
				resources := providers.Resources{
					Provider: string(hostlist.Name),
					Hosts:    []string{"127.0.0.1"},
				}
				resArg, err := functions.AsStarlarkStruct(resources)
				if err != nil {
					t.Fatal(err)
				}
				return []starlark.Tuple{
					{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					{starlark.String("resources"), resArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				res, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatal("unexpected function error: ", result.Error)
				}
			},
		},
		{
			name: "simple cmd",
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"username":         starlark.String(testSupport.CurrentUsername()),
					"port":             starlark.String(testSupport.PortValue()),
					"private_key_path": starlark.String(testSupport.PrivateKeyPath()),
					"max_retries":      starlark.MakeInt(testSupport.MaxConnectionRetries()),
				})
				resArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"provider": starlark.String(hostlist.Name),
					"hosts":    starlark.NewList([]starlark.Value{starlark.String("127.0.0.1")}),
				})
				return []starlark.Tuple{
					{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					{starlark.String("resources"), resArg},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				thread := &starlark.Thread{}
				res, err := Func(thread, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}
				if len(result.Procs) != 1 {
					t.Fatal("missing command result")
				}
				output := strings.TrimSpace(result.Procs[0].Output)
				if output != "Hello World!" {
					t.Error("unexpected result:", output)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var kwargs []starlark.Tuple
			if test.kwargs != nil {
				kwargs = test.kwargs(t)
			}
			test.eval(t, kwargs)
		})
	}
}

func TestRunScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(*testing.T, string)
	}{
		{
			name: "simple command",
			script: fmt.Sprintf(`
result=run(
    cmd="""echo 'Hello World!'""",
    ssh_config = ssh_config(
       username="%s",
        port="%s",
        private_key_path="%s",
        max_retries=%d,
    ),
    resources=hostlist_provider(hosts=["127.0.0.1"]).resources,
)`, testSupport.CurrentUsername(), testSupport.PortValue(), testSupport.PrivateKeyPath(), testSupport.MaxConnectionRetries()),
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}
				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("run() should be assigned to a variable for test")
				}
				var result Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}
				if len(result.Procs) != 1 {
					t.Fatal("missing command result")
				}
				expected := "Hello World!"
				out := strings.TrimSpace(result.Procs[0].Output)
				if out != expected {
					t.Error("unexpected result:", output)
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

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"os"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

func TestScriptConfFunc(t *testing.T) {
	tests := []struct {
		name   string
		kwargs []starlark.Tuple
		eval   func(*testing.T, []starlark.Tuple)
	}{
		{
			name:   "no args",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var res Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Conf.Workdir != DefaultWorkdir() {
					t.Errorf("unexpected workdir value: %s", res.Conf.Workdir)
				}
				if err := os.RemoveAll(res.Conf.Workdir); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:   "with workdir",
			kwargs: []starlark.Tuple{{starlark.String("workdir"), starlark.String("foo")}},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var res Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Conf.Workdir != "foo" {
					t.Errorf("unexpected workdir value: %s", res.Conf.Workdir)
				}
				if err := os.RemoveAll(res.Conf.Workdir); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "with ssh-agent",
			kwargs: []starlark.Tuple{
				{starlark.String("workdir"), starlark.String("foo")},
				{starlark.String("use_ssh_agent"), starlark.Bool(true)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var res Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Conf.Workdir != "foo" {
					t.Errorf("unexpected workdir value: %s", res.Conf.Workdir)
				}
				if err := os.RemoveAll(res.Conf.Workdir); err != nil {
					t.Error(err)
				}
				if !res.Conf.UseSSHAgent {
					t.Errorf("unexpected conf.UseSSHAgent: %t", res.Conf.UseSSHAgent)
				}
			},
		},
		{
			name: "with shell",
			kwargs: []starlark.Tuple{
				{starlark.String("default_shell"), starlark.String("/a/b/c")},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var res Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Conf.DefaultShell != "/a/b/c" {
					t.Errorf("unexpected defaultShell value: %s", res.Conf.DefaultShell)
				}
				if err := os.RemoveAll(res.Conf.Workdir); err != nil {
					t.Error(err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.kwargs)
		})
	}
}

func TestScriptConfScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "run local",
			script: `
result = script_conf(workdir="foo", use_ssh_agent=false)
`,
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}

				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("script_conf() should be assigned to a variable for test")
				}
				var result Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}

				if result.Conf.Workdir != "foo" {
					t.Fatalf("unexpected workdir %s", result.Conf.Workdir)
				}
				if err := os.RemoveAll(result.Conf.Workdir); err != nil {
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

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package make_scriptconf

import (
	"os"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
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
				var res scriptconf.Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Config.Workdir != scriptconf.DefaultWorkdir() {
					t.Errorf("unexpected workdir value: %s", res.Config.Workdir)
				}
				if err := os.RemoveAll(res.Config.Workdir); err != nil {
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
				var res scriptconf.Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Config.Workdir != "foo" {
					t.Errorf("unexpected workdir value: %s", res.Config.Workdir)
				}
				if err := os.RemoveAll(res.Config.Workdir); err != nil {
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
				var res scriptconf.Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Config.Workdir != "foo" {
					t.Errorf("unexpected workdir value: %s", res.Config.Workdir)
				}
				if err := os.RemoveAll(res.Config.Workdir); err != nil {
					t.Error(err)
				}
				if !res.Config.UseSSHAgent {
					t.Errorf("unexpected conf.UseSSHAgent: %t", res.Config.UseSSHAgent)
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
				var res scriptconf.Result
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}
				if res.Config.DefaultShell != "/a/b/c" {
					t.Errorf("unexpected defaultShell value: %s", res.Config.DefaultShell)
				}
				if err := os.RemoveAll(res.Config.Workdir); err != nil {
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

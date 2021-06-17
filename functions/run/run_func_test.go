// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers/hostlist"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
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
				resArg, err := functions.Result(hostlist.Name, providers.Result{Resources: resources})
				if err != nil {
					t.Fatal(err)
				}
				return []starlark.Tuple{
					{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					{starlark.String("resources"), resArg},
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
			name: "simple cmd",
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshConf := sshconf.DefaultConfig()
				sshArg, err := functions.Result(sshconf.Name, sshConf)
				if err != nil {
					t.Fatal(err)
				}
				resources := providers.Resources{
					Provider: string(hostlist.Name),
					Hosts:    []string{"127.0.0.1"},
				}
				resArg, err := functions.Result(hostlist.Name, providers.Result{Resources: resources})
				if err != nil {
					t.Fatal(err)
				}
				return []starlark.Tuple{
					{starlark.String("cmd"), starlark.String("echo 'Hello World!'")},
					{starlark.String("resources"), resArg},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				res, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(result); err != nil {
					t.Fatal(err)
				}
				if len(result.Procs) != 1 {
					t.Error("missing command result")
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

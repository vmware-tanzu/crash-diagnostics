// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package hostlist

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

func TestHostListProviderFunc(t *testing.T) {
	tests := []struct {
		name   string
		kwargs []starlark.Tuple
		eval   func(*testing.T, []starlark.Tuple)
	}{
		{
			name:   "empty args",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err == nil {
					t.Fatal("expecting argument error, got none")
				}
			},
		},
		{
			name: "with hosts",
			kwargs: []starlark.Tuple{
				{
					starlark.String("hosts"),
					starlark.NewList([]starlark.Value{starlark.String("foo"), starlark.String("bar")}),
				},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var result providers.Result
				if err := typekit.Starlark(val).Go(&result); err != nil {
					t.Fatal(err)
				}

				if len(result.Resources.Hosts) != 2 {
					t.Errorf("unexpected host count %d", len(result.Resources.Hosts))
				}
				for i := range result.Resources.Hosts {
					if result.Resources.Hosts[i] != "foo" && result.Resources.Hosts[i] != "bar" {
						t.Errorf("unexpected resource hosts values %s", result.Resources.Hosts[i])
					}
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

func TestHostlistProviderScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(*testing.T, string)
	}{
		{
			name:   "simple script",
			script: `result=hostlist_provider(hosts=["127.0.0.1", "localhost"])`,
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}

				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("hostlist_provider() should be assigned to a variable for testing")
				}
				var result providers.Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}
				if len(result.Resources.Hosts) != 2 {
					t.Errorf("unexpected host count %d", len(result.Resources.Hosts))
				}
				for i := range result.Resources.Hosts {
					if result.Resources.Hosts[i] != "127.0.0.1" && result.Resources.Hosts[i] != "localhost" {
						t.Errorf("unexpected resource hosts values %s", result.Resources.Hosts[i])
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

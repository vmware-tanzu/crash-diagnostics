// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package hostlist

import (
	"testing"

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
				var res providers.Resources
				if err := typekit.Starlark(val).Go(&res); err != nil {
					t.Fatal(err)
				}

				if len(res.Hosts) != 2 {
					t.Errorf("unexpected host count %d", len(res.Hosts))
				}
				for i := range res.Hosts {
					if res.Hosts[i] != "foo" && res.Hosts[i] != "bar" {
						t.Errorf("unexpected resource hosts values %s", res.Hosts[i])
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

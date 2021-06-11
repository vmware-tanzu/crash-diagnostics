// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package hostlist

import (
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"go.starlark.net/starlark"
)

func TestCmd_Run(t *testing.T) {
	tests := []struct {
		name       string
		args       Args
		res        providers.Resources
		shouldFail bool
	}{
		{
			name:       "empty args",
			args:       Args{},
			shouldFail: true,
		},
		{
			name: "multi-hosts",
			args: Args{Hosts: []string{"foo", "bar"}},
			res:  providers.Resources{Hosts: []string{"foo", "bar"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := newCmd()
			args := cmd.Run(new(starlark.Thread), test.args)
			if args.Error != "" && !test.shouldFail {
				t.Fatal("unexpected error:", args.Error)
			}
			if len(args.Hosts) != len(test.args.Hosts) {
				t.Errorf("unexpected host count %d", len(args.Hosts))
			}
			for i := range args.Hosts {
				if args.Hosts[i] != test.args.Hosts[i] {
					t.Errorf("unexpected host %s", args.Hosts[i])
				}
			}
		})
	}
}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestHostListProvider(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name:   "single host",
			script: `provider = host_list_provider(hosts=["foo.host"], ssh_config = ssh_config(username="uname", private_key_path="path"))`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["provider"]
				if data == nil {
					t.Fatalf("%s function not returning value", identifiers.hostListProvider)
				}
				provider, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting *starlark.Struct, got %T", data)
				}
				val, err := provider.Attr("hosts")
				if err != nil {
					t.Fatal(err)
				}
				list := val.(*starlark.List)
				if list.Len() != 1 {
					t.Fatalf("expecting %d items for argument 'hosts', got %d", 2, list.Len())
				}

				sshcfg, err := provider.Attr(identifiers.sshCfg)
				if err != nil {
					t.Error(err)
				}
				if sshcfg == nil {
					t.Errorf("%s missing ssh_config", identifiers.hostListProvider)
				}
			},
		},
		{
			name:   "multiple hosts",
			script: `provider = host_list_provider(hosts=["foo.host.1", "foo.host.2"], ssh_config = ssh_config(username="uname", private_key_path="path"))`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["provider"]
				if data == nil {
					t.Fatalf("%s function not returning value", identifiers.hostListProvider)
				}
				provider, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting *starlark.Struct, got %T", data)
				}

				val, err := provider.Attr("hosts")
				if err != nil {
					t.Fatal(err)
				}
				list := val.(*starlark.List)
				if list.Len() != 2 {
					t.Fatalf("expecting %d items for argument 'hosts', got %d", 2, list.Len())
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

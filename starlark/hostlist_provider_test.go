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
			script: `provider = host_list_provider(hosts="foo.host")`,
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
				if len(provider.AttrNames()) != 1 {
					t.Fatalf("unexpected item count in configs.crashd: %d", len(provider.AttrNames()))
				}
				val, err := provider.Attr("hosts")
				if err != nil {
					t.Fatal(err)
				}
				if trimQuotes(val.String()) != "foo.host" {
					t.Fatalf("unexpected value for key %s in configs.crashd", val.String())
				}
			},
		},
		{
			name:   "multiple hosts",
			script: `provider = host_list_provider(hosts=["foo.host.1", "foo.host.2"])`,
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
				if len(provider.AttrNames()) != 1 {
					t.Fatalf("unexpected item %s: %d", identifiers.hostListProvider, len(provider.AttrNames()))
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

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"strings"
	"testing"
)

func TestInstanceListProvider(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name:   "single instance",
			script: `provider = instance_list_provider(instances=["i-instance"], region="eu-west-1")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["provider"]
				if data == nil {
					t.Fatalf("%s function not returning value", identifiers.instanceListProvider)
				}
				provider, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting *starlark.Struct, got %T", data)
				}
				instances, err := provider.Attr("instances")
				if err != nil {
					t.Fatal(err)
				}
				list := instances.(*starlark.List)
				if list.Len() != 1 {
					t.Fatalf("expecting %d items for argument 'instances', got %d", 1, list.Len())
				}
				region, err := provider.Attr("region")
				if err != nil {
					t.Fatal(err)
				}
				str := region.(starlark.String)
				if str.Len() == 0 {
					t.Fatalf("'region' of '%s' argument cannot be blank", identifiers.instanceListProvider)
				}
			},
		},
		{
			name:   "multiple instance",
			script: `provider = instance_list_provider(instances=["i-instance0", "i-instance1"], region="eu-west-1")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["provider"]
				if data == nil {
					t.Fatalf("%s function not returning value", identifiers.instanceListProvider)
				}
				provider, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting *starlark.Struct, got %T", data)
				}
				instances, err := provider.Attr("instances")
				if err != nil {
					t.Fatal(err)
				}
				list := instances.(*starlark.List)
				if list.Len() != 2 {
					t.Fatalf("expecting %d items for argument 'instances', got %d", 2, list.Len())
				}
				region, err := provider.Attr("region")
				if err != nil {
					t.Fatal(err)
				}
				str := region.(starlark.String)
				if str.Len() == 0 {
					t.Fatalf("'region' of '%s' argument cannot be blank", identifiers.instanceListProvider)
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
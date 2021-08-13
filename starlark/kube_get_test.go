// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"strings"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestKubeGet(t *testing.T) {
	tests := []struct {
		name   string
		kwargs func(t *testing.T) []starlark.Tuple
		eval   func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name: "list of services as starlark objects",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("groups"), starlark.NewList([]starlark.Value{starlark.String("core")})},
					[]starlark.Value{starlark.String("kinds"), starlark.NewList([]starlark.Value{starlark.String("services")})},
					[]starlark.Value{starlark.String("namespaces"), starlark.NewList([]starlark.Value{starlark.String("default"), starlark.String("kube-system")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeGetFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				objVal, err := resultStruct.Attr("objs")
				if err != nil {
					t.Error(err)
				}
				objList, ok := objVal.(*starlark.List)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				if objList.Len() != 2 {
					t.Errorf("unexpected object list returned: %d", objList.Len())
				}
			},
		},
		{
			name: "list of nodes as starlark objects",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("groups"), starlark.NewList([]starlark.Value{starlark.String("core")})},
					[]starlark.Value{starlark.String("kinds"), starlark.NewList([]starlark.Value{starlark.String("nodes")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeGetFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				objVal, err := resultStruct.Attr("objs")
				if err != nil {
					t.Error(err)
				}
				objList, ok := objVal.(*starlark.List)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				if objList.Len() != 1 {
					t.Errorf("unexpected object list returned: %d", objList.Len())
				}
			},
		},
		{
			name: "different categories of objects as starlark objects",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("categories"), starlark.NewList([]starlark.Value{starlark.String("all")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeGetFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				objVal, err := resultStruct.Attr("objs")
				if err != nil {
					t.Error(err)
				}
				objList, ok := objVal.(*starlark.List)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				if objList.Len() <= 1 {
					t.Errorf("unexpected object list returned: %d", objList.Len())
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.kwargs(t))
		})
	}
}

func TestKubeGetScript(t *testing.T) {
	k8sconfig := testSupport.KindKubeConfigFile()
	clusterName := testSupport.KindClusterContextName()

	execute := func(t *testing.T, script string) *starlarkstruct.Struct {
		executor := New()
		if err := executor.Exec("test.kube.capture", strings.NewReader(script)); err != nil {
			t.Fatalf("failed to exec: %s", err)
		}
		if !executor.result.Has("kube_data") {
			t.Fatalf("script result must be assigned to a value")
		}

		data, ok := executor.result["kube_data"].(*starlarkstruct.Struct)
		if !ok {
			t.Fatal("script result is not a struct")
		}
		return data
	}

	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "namespaced objects as starlark objects with context",
			script: fmt.Sprintf(`
set_defaults(kube_config(path="%s", cluster_context="%s"))
kube_data = kube_get(groups=["core"], kinds=["services"], namespaces=["default", "kube-system"])
`, k8sconfig, clusterName),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				errVal, err := data.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				objVal, err := data.Attr("objs")
				if err != nil {
					t.Error(err)
				}
				objList, ok := objVal.(*starlark.List)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				if objList.Len() != 2 {
					t.Errorf("unexpected object list returned: %d", objList.Len())
				}
			},
		},
		{
			name: "non-namespaced objects as starlark objects",
			script: fmt.Sprintf(`
set_defaults(kube_config(path="%s"))
kube_data = kube_get(groups=["core"], kinds=["nodes"])
`, k8sconfig),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				errVal, err := data.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				objVal, err := data.Attr("objs")
				if err != nil {
					t.Error(err)
				}
				objList, ok := objVal.(*starlark.List)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				if objList.Len() != 1 {
					t.Errorf("unexpected object list returned: %d", objList.Len())
				}
			},
		},
		{
			name: "different categories of objects as starlark objects with context",
			script: fmt.Sprintf(`
set_defaults(kube_config(path="%s", cluster_context="%s"))
kube_data = kube_get(categories=["all"])
`, k8sconfig, clusterName),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				errVal, err := data.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				objVal, err := data.Attr("objs")
				if err != nil {
					t.Error(err)
				}
				objList, ok := objVal.(*starlark.List)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				if objList.Len() < 3 {
					t.Errorf("unexpected object list returned: %d", objList.Len())
				}
			},
		},
		{
			name: "retrieve containers as starlark objects",
			script: fmt.Sprintf(`
set_defaults(kube_config(path="%s"))
kube_data = kube_get(kinds=["pods"], namespaces=["kube-system"], containers=["etcd"])
`, k8sconfig),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				errVal, err := data.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				objVal, err := data.Attr("objs")
				if err != nil {
					t.Error(err)
				}
				objList, ok := objVal.(*starlark.List)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				if objList.Len() < 1 {
					t.Errorf("unexpected object list returned: %d", objList.Len())
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

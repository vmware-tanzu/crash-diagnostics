// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestResourcesFunc(t *testing.T) {
	tests := []struct {
		name   string
		kwargs func(t *testing.T) []starlark.Tuple
		eval   func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name:   "empty kwargs",
			kwargs: func(t *testing.T) []starlark.Tuple { return nil },
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := resourcesFunc(&starlark.Thread{Name: "test"}, nil, nil, kwargs)
				if err == nil {
					t.Fatal("expected failure, but err == nil")
				}
			},
		},
		{
			name: "bad args",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{[]starlark.Value{starlark.String("foo"), starlark.String("bar")}}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := resourcesFunc(&starlark.Thread{Name: "test"}, nil, nil, kwargs)
				if err == nil {
					t.Fatal("expected failure, but err == nil")
				}
			},
		},
		{
			name: "missing ssh_config",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{[]starlark.Value{starlark.String("hosts"), starlark.String("foo.host.1")}}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := resourcesFunc(&starlark.Thread{Name: "test"}, nil, nil, kwargs)
				if err == nil {
					t.Fatal("expected failure, but err == nil")
				}
			},
		},
		{
			name: "host only",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("hosts"), starlark.String("foo.host.1")},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				res, err := resourcesFunc(newThreadLocal(), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				resStruct, ok := res.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for resource: %T", res)
				}
				val, err := resStruct.Attr("kind")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(val.String()) != identifiers.hostListResources {
					t.Errorf("unexpected resource kind for host list provider")
				}

				transport, err := resStruct.Attr("transport")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(transport.String()) != "ssh" {
					t.Errorf("unexpected %s transport: %s", identifiers.resources, transport)
				}

				sshCfg, err := resStruct.Attr(identifiers.sshCfg)
				if err != nil {
					t.Error(err)
				}
				if sshCfg == nil {
					t.Error("resources missing ssh_config")
				}

				hosts, err := resStruct.Attr("hosts")
				if err != nil {
					t.Error(err)
				}
				hostList := hosts.(*starlark.List)
				if trimQuotes(hostList.Index(0).String()) != "foo.host.1" {
					t.Error("unexpected value for names list in resources")
				}
			},
		},
		{
			name: "provider only",
			kwargs: func(t *testing.T) []starlark.Tuple {
				provider, err := newHostListProvider(
					newThreadLocal(),
					starlark.StringDict{"hosts": starlark.NewList([]starlark.Value{starlark.String("local.host")})},
				)
				if err != nil {
					t.Fatal(err)
				}

				return []starlark.Tuple{[]starlark.Value{starlark.String("provider"), provider}}
			},

			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				res, err := resourcesFunc(newThreadLocal(), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				resStruct, ok := res.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for resource: %T", res)
				}
				val, err := resStruct.Attr("kind")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(val.String()) != identifiers.hostListResources {
					t.Errorf("unexpected resource kind for host list provider")
				}

				transport, err := resStruct.Attr("transport")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(transport.String()) != "ssh" {
					t.Errorf("unexpected %s transport: %s", identifiers.resources, transport)
				}

				sshCfg, err := resStruct.Attr(identifiers.sshCfg)
				if err != nil {
					t.Error(err)
				}
				if sshCfg == nil {
					t.Error("resources missing ssh_config")
				}

				hosts, err := resStruct.Attr("hosts")
				if err != nil {
					t.Error(err)
				}
				hostList := hosts.(*starlark.List)
				if trimQuotes(hostList.Index(0).String()) != "local.host" {
					t.Error("unexpected value for names list in resources")
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

func TestResourceScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name:   "default resource with host",
			script: `resources(hosts="foo.host.1")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.thread.Local(identifiers.resources)
				if data == nil {
					t.Fatalf("default %s not found in thread", identifiers.resources)
				}
				resStruct, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting *starlark.Struct, got %T", data)
				}

				val, err := resStruct.Attr("kind")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(val.String()) != identifiers.hostListResources {
					t.Errorf("unexpected resource kind for host list provider")
				}

				transport, err := resStruct.Attr("transport")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(transport.String()) != "ssh" {
					t.Errorf("unexpected %s transport: %s", identifiers.resources, transport)
				}

				sshCfg, err := resStruct.Attr(identifiers.sshCfg)
				if err != nil {
					t.Error(err)
				}
				if sshCfg == nil {
					t.Error("resources missing ssh_config")
				}

				hosts, err := resStruct.Attr("hosts")
				if err != nil {
					t.Error(err)
				}
				hostList := hosts.(*starlark.List)
				if trimQuotes(hostList.Index(0).String()) != "foo.host.1" {
					t.Error("unexpected value for names list in resources")
				}
			},
		},
		{
			name:   "default resource with provider",
			script: `resources(provider=host_list_provider(hosts="foo.host.1"))`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.thread.Local(identifiers.resources)
				if data == nil {
					t.Fatalf("default %s not found in thread", identifiers.resources)
				}
				resStruct, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting *starlark.Struct, got %T", data)
				}

				val, err := resStruct.Attr("kind")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(val.String()) != identifiers.hostListResources {
					t.Errorf("unexpected resource kind for host list provider")
				}

				transport, err := resStruct.Attr("transport")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(transport.String()) != "ssh" {
					t.Errorf("unexpected %s transport: %s", identifiers.resources, transport)
				}

				sshCfg, err := resStruct.Attr(identifiers.sshCfg)
				if err != nil {
					t.Error(err)
				}
				if sshCfg == nil {
					t.Error("resources missing ssh_config")
				}

				hosts, err := resStruct.Attr("hosts")
				if err != nil {
					t.Error(err)
				}
				hostList := hosts.(*starlark.List)
				if trimQuotes(hostList.Index(0).String()) != "foo.host.1" {
					t.Error("unexpected value for names list in resources")
				}
			},
		},
		{
			name:   "resources assigned",
			script: `res = resources(hosts="foo.host.1")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["res"]
				if data == nil {
					t.Fatalf("%s function call not returning value", identifiers.resources)
				}
				resStruct, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting *starlark.Struct, got %T", data)
				}

				val, err := resStruct.Attr("kind")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(val.String()) != identifiers.hostListResources {
					t.Errorf("unexpected resource kind for host list provider")
				}

				transport, err := resStruct.Attr("transport")
				if err != nil {
					t.Error(err)
				}
				if trimQuotes(transport.String()) != "ssh" {
					t.Errorf("unexpected %s transport: %s", identifiers.resources, transport)
				}

				sshCfg, err := resStruct.Attr(identifiers.sshCfg)
				if err != nil {
					t.Error(err)
				}
				if sshCfg == nil {
					t.Error("resources missing ssh_config")
				}

				hosts, err := resStruct.Attr("hosts")
				if err != nil {
					t.Error(err)
				}
				hostList := hosts.(*starlark.List)
				if trimQuotes(hostList.Index(0).String()) != "foo.host.1" {
					t.Error("unexpected value for names list in resources")
				}
			},
		},
		//{
		//	name:   "multiple hosts",
		//	script: `provider = host_list_provider(hosts=["foo.host.1", "foo.host.2"])`,
		//	eval: func(t *testing.T, script string) {
		//		exe := New()
		//		if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
		//			t.Fatal(err)
		//		}
		//		data := exe.result["provider"]
		//		if data == nil {
		//			t.Fatalf("%s function not returning value", identifiers.hostListProvider)
		//		}
		//		provider, ok := data.(*starlarkstruct.Struct)
		//		if !ok {
		//			t.Fatalf("expecting *starlark.Struct, got %T", data)
		//		}
		//
		//		val, err := provider.Attr("hosts")
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		list := val.(*starlark.List)
		//		if list.Len() != 2 {
		//			t.Fatalf("expecting %d items for argument 'hosts', got %d", 2, list.Len())
		//		}
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script)
		})
	}
}

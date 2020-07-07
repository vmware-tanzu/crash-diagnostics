// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"
	"testing"

	"go.starlark.net/starlarkstruct"
)

func TestExec(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name:   "crash_config only",
			script: `crashd_config()`,
			eval: func(t *testing.T, script string) {
				if err := New().Exec("test.file", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:   "kube_config only",
			script: `kube_config()`,
			eval: func(t *testing.T, script string) {
				if err := New().Exec("test.file", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:   "kube_config only",
			script: `kube_config()`,
			eval: func(t *testing.T, script string) {
				if err := New().Exec("test.file", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:   "args only",
			script: `cfg = ssh_config(username=args.username)`,
			eval: func(t *testing.T, script string) {
				e := New()
				e.WithArgs(map[string]string{
					"username": "foo",
				})
				if err := e.Exec("test.file", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := e.result["cfg"]
				if data == nil {
					t.Fatal("ssh_config function not returning value")
				}
				cfg, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for thread local key ssh_config: %T", data)
				}

				val, err := cfg.Attr("username")
				if err != nil {
					t.Fatal(err)
				}
				if trimQuotes(val.String()) != "foo" {
					t.Fatalf("unexpected value for key 'username': %s", val.String())
				}
			},
		},
		{
			name:   "multiple args",
			script: `cfg = ssh_config(username=args.username, port=args.ssh_port)`,
			eval: func(t *testing.T, script string) {
				e := New()
				e.WithArgs(map[string]string{
					"username": "bar",
					"ssh_port": "1234",
				})
				if err := e.Exec("test.file", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := e.result["cfg"]
				if data == nil {
					t.Fatal("ssh_config function not returning value")
				}
				cfg, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for thread local key ssh_config: %T", data)
				}

				val, err := cfg.Attr("username")
				if err != nil {
					t.Fatal(err)
				}
				if trimQuotes(val.String()) != "bar" {
					t.Fatalf("unexpected value for key 'username': %s", val.String())
				}

				portVal, err := cfg.Attr("port")
				if err != nil {
					t.Fatal(err)
				}
				if trimQuotes(portVal.String()) != "1234" {
					t.Fatalf("unexpected value for key 'port': %s", val.String())
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

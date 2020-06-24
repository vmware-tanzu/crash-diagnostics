// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"
	"testing"

	"go.starlark.net/starlarkstruct"
)

func TestSSHConfigNew(t *testing.T) {
	e := New()
	if e.thread == nil {
		t.Error("thread is nil")
	}
	cfg := e.thread.Local(identifiers.sshCfg)
	if cfg == nil {
		t.Error("ssh_config dict not found in thread")
	}
}

func TestSSHConfigFunc(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name:   "ssh_config saved in thread",
			script: `ssh_config(username="uname", private_key_path="path")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.thread.Local(identifiers.sshCfg)
				if data == nil {
					t.Fatal("ssh_config not saved in thread local")
				}
				cfg, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for thread local key ssh_config: %T", data)
				}
				if len(cfg.AttrNames()) != 4 {
					t.Fatalf("unexpected item count in ssh_config: %d", len(cfg.AttrNames()))
				}
				val, err := cfg.Attr("username")
				if err != nil {
					t.Fatal(err)
				}
				if trimQuotes(val.String()) != "uname" {
					t.Fatalf("unexpected value for key 'foo': %s", val.String())
				}
			},
		},

		{
			name:   "ssh_config returned value",
			script: `cfg = ssh_config(username="uname", private_key_path="path")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["cfg"]
				if data == nil {
					t.Fatal("ssh_config function not returning value")
				}
				cfg, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for thread local key ssh_config: %T", data)
				}
				if len(cfg.AttrNames()) != 4 {
					t.Fatalf("unexpected item count in ssh_config: %d", len(cfg.AttrNames()))
				}
				val, err := cfg.Attr("private_key_path")
				if err != nil {
					t.Fatal(err)
				}
				if trimQuotes(val.String()) != "path" {
					t.Fatalf("unexpected value for key %s in ssh_config", val.String())
				}
			},
		},

		{
			name:   "crash_config default",
			script: `one = 1`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.thread.Local(identifiers.sshCfg)
				if data == nil {
					t.Fatal("default ssh_config not saved in thread local")
				}

				cfg, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for thread local key ssh_config: %T", data)
				}
				if len(cfg.AttrNames()) != 5 {
					t.Fatalf("unexpected item count in ssh_config: %d", len(cfg.AttrNames()))
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

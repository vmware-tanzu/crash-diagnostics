// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestCrashdConfigNew(t *testing.T) {
	e := New()
	if e.thread == nil {
		t.Error("thread is nil")
	}
	cfg := e.thread.Local(identifiers.crashdCfg)
	if cfg == nil {
		t.Error("crashd_config dict not found in thread")
	}
}

func TestCrashdConfigFunc(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name:   "crash_config saved in thread",
			script: `crashd_config(foo="fooval", bar="barval")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.thread.Local(identifiers.crashdCfg)
				if data == nil {
					t.Fatal("crashd_config not saved in thread local")
				}
				cfg, ok := data.(*starlark.Dict)
				if !ok {
					t.Fatalf("unexpected type for thread local key configs.crashd: %T", data)
				}
				if cfg.Len() != 2 {
					t.Fatalf("unexpected item count in configs.crashd: %d", cfg.Len())
				}
				val, found, err := cfg.Get(starlark.String("foo"))
				if err != nil {
					t.Fatal(err)
				}
				if !found {
					t.Fatalf("key 'foo' not found in configs.crashd")
				}
				if trimQuotes(val.String()) != "fooval" {
					t.Fatalf("unexpected value for key 'foo': %s", val.String())
				}
			},
		},

		{
			name:   "crash_config returned value",
			script: `cfg = crashd_config(foo="fooval", bar="barval")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["cfg"]
				if data == nil {
					t.Fatal("crashd_config function not returning value")
				}
				cfg, ok := data.(*starlark.Dict)
				if !ok {
					t.Fatalf("unexpected type for thread local key configs.crashd: %T", data)
				}
				if cfg.Len() != 2 {
					t.Fatalf("unexpected item count in configs.crashd: %d", cfg.Len())
				}
				val, found, err := cfg.Get(starlark.String("foo"))
				if err != nil {
					t.Fatal(err)
				}
				if !found {
					t.Fatalf("key 'foo' not found in configs.crashd")
				}
				if trimQuotes(val.String()) != "fooval" {
					t.Fatalf("unexpected value for key %s in configs.crashd", val.String())
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
				data := exe.thread.Local(identifiers.crashdCfg)
				if data == nil {
					t.Fatal("default crashd_config not saved in thread local")
				}

				cfg, ok := data.(*starlark.Dict)
				if !ok {
					t.Fatalf("unexpected type for thread local key crashd_config: %T", data)
				}
				if cfg.Len() != 4 {
					t.Fatalf("unexpected item count in configs.crashd: %d", cfg.Len())
				}
				val, found, err := cfg.Get(starlark.String("uid"))
				if err != nil {
					t.Fatal(err)
				}
				if !found {
					t.Fatalf("key 'foo' not found in configs.crashd")
				}
				if trimQuotes(val.String()) != getUid() {
					t.Fatalf("unexpected value for key %s in configs.crashd", val.String())
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

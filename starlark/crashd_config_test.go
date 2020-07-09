// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"os"
	"strings"
	"testing"

	"go.starlark.net/starlarkstruct"
)

func testCrashdConfigNew(t *testing.T) {
	e := New()
	if e.thread == nil {
		t.Error("thread is nil")
	}
}

func testCrashdConfigFunc(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name:   "crash_config saved in thread",
			script: `crashd_config(workdir="fooval", default_shell="barval")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.thread.Local(identifiers.crashdCfg)
				if data == nil {
					t.Fatal("crashd_config not saved in thread local")
				}
				cfg, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for thread local key configs.crashd: %T", data)
				}
				if len(cfg.AttrNames()) != 5 {
					t.Fatalf("unexpected item count in configs.crashd: %d", len(cfg.AttrNames()))
				}

			},
		},

		{
			name:   "crash_config returned value",
			script: `cfg = crashd_config(uid="fooval", gid="barval")`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				data := exe.result["cfg"]
				if data == nil {
					t.Fatal("crashd_config function not returning value")
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

				cfg, ok := data.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type for thread local key crashd_config: %T", data)
				}
				if len(cfg.AttrNames()) != 5 {
					t.Fatalf("unexpected item count in configs.crashd: %d", len(cfg.AttrNames()))
				}
				val, err := cfg.Attr("uid")
				if err != nil {
					t.Fatalf("key 'foo' not found in configs.crashd: %s", err)
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

func TestCrashdCfgAll(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{name: "testCrashdConfigNew", test: testCrashdConfigNew},
		{name: "testCrashdConfigFunc", test: testCrashdConfigFunc},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer os.RemoveAll(defaults.workdir)
			test.test(t)
		})
	}
}

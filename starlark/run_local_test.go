// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestRunLocalFunc(t *testing.T) {
	tests := []struct {
		name string
		args func(t *testing.T) starlark.Tuple
		eval func(t *testing.T, args starlark.Tuple)
	}{
		{
			name: "simple command",
			args: func(t *testing.T) starlark.Tuple { return starlark.Tuple{starlark.String("echo 'Hello World!'")} },
			eval: func(t *testing.T, args starlark.Tuple) {
				val, err := runLocalFunc(newTestThreadLocal(t), nil, args, nil)
				if err != nil {
					t.Fatal(err)
				}
				result := ""
				if r, ok := val.(starlark.String); ok {
					result = string(r)
				}
				if result != "Hello World!" {
					t.Errorf("unexpected result: %s", result)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.args(t))
		})
	}
}

func TestRunLocalScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "run local",
			script: `
result = run_local("""echo 'Hello World!'""")
`,
			eval: func(t *testing.T, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("run_local() should be assigned to a variable for test")
				}
				result, ok := resultVal.(starlark.String)
				if !ok {
					t.Fatal("run_local() should return a string")
				}

				if string(result) != "Hello World!" {
					t.Fatalf("uneexpected result %s", result)
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

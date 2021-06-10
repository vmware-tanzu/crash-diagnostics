// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run_local

import (
	"strings"
	"testing"

	crashlark "github.com/vmware-tanzu/crash-diagnostics/starlark"
	crashtest "github.com/vmware-tanzu/crash-diagnostics/testing"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"

	"go.starlark.net/starlark"
)

func TestRunLocalFunc(t *testing.T) {
	tests := []struct {
		name   string
		kwargs []starlark.Tuple
		eval   func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name:   "simple command",
			kwargs: []starlark.Tuple{{starlark.String("cmd"), starlark.String("echo 'Hello World!'")}},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := runLocalFunc(crashtest.NewStarlarkThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var p Result
				if err := typekit.Starlark(val).Go(&p); err != nil {
					t.Fatalf("unable to convert result: %s", err)
				}

				if p.Result != "Hello World!" {
					t.Errorf("unexpected result: %s", p.Result)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.kwargs)
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
				exe := crashlark.New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.Result()["result"]
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

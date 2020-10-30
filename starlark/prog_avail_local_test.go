// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestProgAvailLocalScript(t *testing.T) {
	tests := []struct {
		name          string
		scriptSnippet string
		exists        bool
	}{
		{
			name:          "prog_avail_local checks for 'go' using positional args",
			scriptSnippet: "prog_avail_local(prog='go')",
			exists:        true,
		},
		{
			name:          "prog_avail_local checks for 'go' using keyword args",
			scriptSnippet: "prog_avail_local('go')",
			exists:        true,
		},
		{
			name:          "prog_avail_local checks for 'nonexistant' using positional args",
			scriptSnippet: "prog_avail_local(prog='nonexistant')",
			exists:        false,
		},
		{
			name:          "prog_avail_local checks for 'nonexistant' using keyword args",
			scriptSnippet: "prog_avail_local('nonexistant')",
			exists:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script := fmt.Sprintf(`path=%v`, test.scriptSnippet)
			exe := New()
			if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
				t.Fatal(err)
			}

			resultVal := exe.result["path"]
			if resultVal == nil {
				t.Fatal("prog_avail_local() should be assigned to a variable for test")
			}

			result, ok := resultVal.(starlark.String)
			if !ok {
				t.Fatal("prog_avail_local() should return a string")
			}

			if (len(string(result)) == 0) == test.exists {
				if test.exists {
					t.Fatalf("expecting prog to exists but 'prog_avail_local' didnt find it")
				} else {
					t.Fatalf("expecting prog to not exists but 'prog_avail_local' found it")

				}
			}
		})
	}
}

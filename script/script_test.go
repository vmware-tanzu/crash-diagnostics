// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"testing"
)

func TestScriptAddDirective(t *testing.T) {
	tests := []struct {
		name   string
		script func(t *testing.T) *Script
		eval   func(*testing.T, *Script)
	}{
		{
			name: "Script.AddDirective/all params",
			script: func(t *testing.T) *Script {
				return New().AddDirective(0, "FOO", `a:"bc" d:"ef"`, ArgMap{"a": "bc", "d": "ef"})
			},
			eval: func(t *testing.T, scr *Script) {
				if len(scr.directives) != 1 {
					t.Error("directive not added")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script(t))
		})
	}
}

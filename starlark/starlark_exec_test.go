// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"
	"testing"
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script)
		})
	}
}

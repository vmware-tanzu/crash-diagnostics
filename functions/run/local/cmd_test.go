// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"testing"

	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/functions/run"
)

func TestCmd_Run(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		expected   string
		shouldFail bool
	}{
		{
			name:     "simple exec",
			param:    `echo "Hello World!"`,
			expected: "Hello World!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := newCmd().Run(&starlark.Thread{}, test.param)
			if !test.shouldFail && err != nil {
				t.Fatal(err)
			}

			proc, ok := result.Value().(run.LocalProc)
			if !ok {
				t.Fatalf("unexpected type in CommandResult %T", result.Value())
			}
			if proc.Result != test.expected {
				t.Errorf("command returned unexpected result: %s", proc.Result)
			}
			if proc.Pid == 0 {
				t.Errorf("successful command returned 0 pid")
			}
			if proc.ExitCode != 0 {
				t.Errorf("successful command returned non-zero exit code")
			}
		})
	}
}

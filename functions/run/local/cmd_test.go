// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestCmd_Run(t *testing.T) {
	tests := []struct {
		name       string
		args       Args
		expected   string
		shouldFail bool
	}{
		{
			name:     "simple exec",
			args:     Args{Cmd: `echo "Hello World!"`},
			expected: "Hello World!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := newCmd().Run(&starlark.Thread{}, test.args)
			if !test.shouldFail && result.Error != "" {
				t.Fatal(result.Error)
			}

			if result.Result != test.expected {
				t.Errorf("command returned unexpected result: %s", result.Result)
			}
			if result.Pid == 0 {
				t.Errorf("successful command returned 0 pid")
			}
			if result.ExitCode != 0 {
				t.Errorf("successful command returned non-zero exit code")
			}
		})
	}
}

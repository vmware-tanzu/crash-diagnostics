// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecFROM(t *testing.T) {
	tests := []execTest{
		{
			name: "FROM with multiple addresses",
			source: func() string {
				return `
				ENV host=local
				FROM '$host'
				`
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

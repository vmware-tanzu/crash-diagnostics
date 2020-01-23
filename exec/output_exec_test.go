// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecOUTPUT(t *testing.T) {
	tests := []execTest{
		{
			name: "exec with OUTPUT",
			source: func() string {
				return fmt.Sprintf("FROM 127.0.0.1:%s\nOUTPUT path:/tmp/crashout/out.tar.gz\nCAPTURE /bin/echo HELLO", testSSHPort)
			},
			exec: func(s *script.Script) error {
				output := s.Preambles[script.CmdOutput][0].(*script.OutputCommand)
				defer os.RemoveAll(output.Path())

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				if _, err := os.Stat(output.Path()); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "exec OUTPUT with var expansion",
			source: func() string {
				return fmt.Sprintf(`
				FROM 127.0.0.1:%s
				ENV outfile=out.tar.gz
				CAPTURE /bin/echo HELLO
				OUTPUT path:/tmp/crashout/${outfile}
				`, testSSHPort)
			},
			exec: func(s *script.Script) error {
				output := s.Preambles[script.CmdOutput][0].(*script.OutputCommand)
				defer os.RemoveAll(output.Path())

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				if _, err := os.Stat(output.Path()); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "exec with missing OUTPUT",
			source: func() string {
				return fmt.Sprintf("FROM 127.0.0.1:%s\nCAPTURE /bin/echo HELLO", testSSHPort)
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

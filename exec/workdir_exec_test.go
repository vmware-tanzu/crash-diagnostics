// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecWORKDIR(t *testing.T) {
	tests := []execTest{
		{
			name: "exec with WORKDIR",
			source: func() string {
				return "WORKDIR /tmp/foodir\nCAPTURE /bin/echo HELLO"
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Host()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				defer os.RemoveAll(workdir.Path())
				capCmd := s.Actions[0].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				fileName := filepath.Join(workdir.Path(), machine, fmt.Sprintf("%s.txt", sanitizeStr(capCmd.GetCmdString())))
				if _, err := os.Stat(fileName); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "exec with default WORKDIR",
			source: func() string {
				return "CAPTURE /bin/echo HELLO"
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Host()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				capCmd := s.Actions[0].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				fileName := filepath.Join(workdir.Path(), machine, fmt.Sprintf("%s.txt", sanitizeStr(capCmd.GetCmdString())))
				if _, err := os.Stat(fileName); err != nil {
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

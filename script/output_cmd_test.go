// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"testing"
)

func TestCommandOUTPUT(t *testing.T) {
	tests := []commandTest{
		{
			name: "OUTPUT without named arg",
			source: func() string {
				return "OUTPUT foo/bar.tar.gz"
			},
			script: func(s *Script) error {
				outs := s.Preambles[CmdOutput]
				if len(outs) != 1 {
					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
				}
				outCmd, ok := outs[0].(*OutputCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", outs[0])
				}
				if outCmd.Path() != "foo/bar.tar.gz" {
					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}
				return nil
			},
		},
		{
			name: "OUTPUT with single arg",
			source: func() string {
				return "OUTPUT path:foo/bar.tar.gz"
			},
			script: func(s *Script) error {
				outs := s.Preambles[CmdOutput]
				if len(outs) != 1 {
					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
				}
				outCmd, ok := outs[0].(*OutputCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", outs[0])
				}
				if outCmd.Path() != "foo/bar.tar.gz" {
					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}
				return nil
			},
		},
		{
			name: "Multiple OUTPUTs",
			source: func() string {
				return "OUTPUT path:foo/bar\nOUTPUT bazz/buzz.tar.gz"
			},
			script: func(s *Script) error {
				outs := s.Preambles[CmdOutput]
				if len(outs) != 1 {
					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
				}
				outCmd, ok := outs[0].(*OutputCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", outs[0])
				}
				if outCmd.Path() != "bazz/buzz.tar.gz" {
					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}
				return nil
			},
		},
		{
			name: "OUTPUT with multiple args",
			source: func() string {
				return "OUTPUT path:foo/bar path:bazz/buzz"
			},
			shouldFail: true,
		},
		{
			name: "OUTPUT with no args",
			source: func() string {
				return "OUTPUT"
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestCommandOUTPUT(t *testing.T) {
	tests := []parserTest{
		{
			name: "OUTPUT",
			source: func(t *testing.T) string {
				return "OUTPUT foo/bar.tar.gz"
			},
			script: func(t *testing.T, s *script.Script) {
				outs := s.Preambles[script.CmdOutput]
				if len(outs) != 1 {
					t.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
				}
				outCmd, ok := outs[0].(*script.OutputCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", outs[0])
				}
				if outCmd.Path() != "foo/bar.tar.gz" {
					t.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}
			},
		},
		//		{
		//			name: "OUTPUT with quoted default param",
		//			source: func() string {
		//				return "OUTPUT 'foo/bar.tar.gz'"
		//			},
		//			script: func(s *Script) error {
		//				outs := s.Preambles[CmdOutput]
		//				if len(outs) != 1 {
		//					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
		//				}
		//				outCmd, ok := outs[0].(*OutputCommand)
		//				if !ok {
		//					return fmt.Errorf("Unexpected type %T in script", outs[0])
		//				}
		//				if outCmd.Path() != "foo/bar.tar.gz" {
		//					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "OUTPUT with single arg",
		//			source: func() string {
		//				return "OUTPUT path:foo/bar.tar.gz"
		//			},
		//			script: func(s *Script) error {
		//				outs := s.Preambles[CmdOutput]
		//				if len(outs) != 1 {
		//					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
		//				}
		//				outCmd, ok := outs[0].(*OutputCommand)
		//				if !ok {
		//					return fmt.Errorf("Unexpected type %T in script", outs[0])
		//				}
		//				if outCmd.Path() != "foo/bar.tar.gz" {
		//					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "Multiple OUTPUTs",
		//			source: func() string {
		//				return "OUTPUT path:foo/bar\nOUTPUT path:'bazz/buzz.tar.gz'"
		//			},
		//			script: func(s *Script) error {
		//				outs := s.Preambles[CmdOutput]
		//				if len(outs) != 1 {
		//					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
		//				}
		//				outCmd, ok := outs[0].(*OutputCommand)
		//				if !ok {
		//					return fmt.Errorf("Unexpected type %T in script", outs[0])
		//				}
		//				if outCmd.Path() != "bazz/buzz.tar.gz" {
		//					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "OUTPUT with expanded var",
		//			source: func() string {
		//				os.Setenv("foopath", "foo/bar.tar.gz")
		//				return "OUTPUT $foopath"
		//			},
		//			script: func(s *Script) error {
		//				outs := s.Preambles[CmdOutput]
		//				if len(outs) != 1 {
		//					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
		//				}
		//				outCmd, ok := outs[0].(*OutputCommand)
		//				if !ok {
		//					return fmt.Errorf("Unexpected type %T in script", outs[0])
		//				}
		//				if outCmd.Path() != "foo/bar.tar.gz" {
		//					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "OUTPUT with multiple args",
		//			source: func() string {
		//				return "OUTPUT path:foo/bar path:bazz/buzz"
		//			},
		//			shouldFail: true,
		//		},
		//		{
		//			name: "OUTPUT with no args",
		//			source: func() string {
		//				return "OUTPUT"
		//			},
		//			shouldFail: true,
		//		},
		//		{
		//			name: "OUTPUT named arg with embedded colon",
		//			source: func() string {
		//				return "OUTPUT path:foo/bar.tar.gz:ignore"
		//			},
		//			script: func(s *Script) error {
		//				outs := s.Preambles[CmdOutput]
		//				if len(outs) != 1 {
		//					return fmt.Errorf("Script has unexpected number of OUTPUT %d", len(outs))
		//				}
		//				outCmd, ok := outs[0].(*OutputCommand)
		//				if !ok {
		//					return fmt.Errorf("Unexpected type %T in script", outs[0])
		//				}
		//				if outCmd.Path() != "foo/bar.tar.gz:ignore" {
		//					return fmt.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
		//				}
		//				return nil
		//			},
		//		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runParserTest(t, test)
		})
	}
}

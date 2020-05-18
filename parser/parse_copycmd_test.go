// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package parser
//
//import (
//	"fmt"
//	"os"
//	"testing"
//)
//
//func TestCommandCOPY(t *testing.T) {
//	tests := []commandTest{
//		{
//			name: "COPY with default param",
//			source: func() string {
//				return "COPY /a/b/c"
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 1 {
//					return fmt.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
//				}
//
//				arg := cmd.Paths()[0]
//				if arg != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument %s", arg)
//				}
//				return nil
//			},
//		},
//		{
//			name: "COPY with quoted default param",
//			source: func() string {
//				return `COPY '/a/b/c'`
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 1 {
//					return fmt.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
//				}
//
//				arg := cmd.Paths()[0]
//				if arg != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument %s", arg)
//				}
//				return nil
//			},
//		},
//		{
//			name: "COPY with quoted named param",
//			source: func() string {
//				return `COPY paths:"/a/b/c"`
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 1 {
//					return fmt.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
//				}
//
//				arg := cmd.Paths()[0]
//				if arg != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument %s", arg)
//				}
//				return nil
//			},
//		},
//		{
//			name: "COPY with multiple args",
//			source: func() string {
//				return "COPY /a/b/c /e/f/g"
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 2 {
//					return fmt.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
//				}
//				if cmd.Paths()[0] != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
//				}
//				if cmd.Paths()[1] != "/e/f/g" {
//					return fmt.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
//				}
//
//				return nil
//			},
//		},
//		{
//			name: "Multiple COPY commands",
//			source: func() string {
//				return "COPY /a/b/c\nCOPY d /e/f"
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 2 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd0 := s.Actions[0].(*CopyCommand)
//				if len(cmd0.Paths()) != 1 {
//					return fmt.Errorf("COPY action[0] has wrong number of args %s", cmd0.Paths())
//				}
//				arg := cmd0.Paths()[0]
//				if arg != "/a/b/c" {
//					return fmt.Errorf("COPY action[0] has unexpected arg %s", arg)
//				}
//
//				cmd1 := s.Actions[1].(*CopyCommand)
//				if len(cmd1.Paths()) != 2 {
//					return fmt.Errorf("COPY action[1] has wrong number of args %d", len(cmd1.Paths()))
//				}
//				arg = cmd1.Paths()[0]
//				if arg != "d" {
//					return fmt.Errorf("COPY action[1] has unexpected arg[0] %s", arg)
//				}
//				arg = cmd1.Paths()[1]
//				if arg != "/e/f" {
//					return fmt.Errorf("COPY action[1] has unexpected arg[1] %s", arg)
//				}
//				return nil
//			},
//		},
//		{
//			name: "COPY single with named param",
//			source: func() string {
//				return "COPY paths:/a/b/c"
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 1 {
//					return fmt.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
//				}
//
//				arg := cmd.Paths()[0]
//				if arg != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument %s", arg)
//				}
//				return nil
//			},
//		},
//		{
//			name: "COPY multiple with named param",
//			source: func() string {
//				return `COPY paths:"/a/b/c /e/f/g"`
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 2 {
//					return fmt.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
//				}
//				if cmd.Paths()[0] != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
//				}
//				if cmd.Paths()[1] != "/e/f/g" {
//					return fmt.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
//				}
//
//				return nil
//			},
//		},
//		{
//			name: "COPY with var expansion",
//			source: func() string {
//				os.Setenv("foopath1", "/a/b/c")
//				os.Setenv("foodir", "g")
//				return "COPY ${foopath1} /e/f/${foodir}"
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 2 {
//					return fmt.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
//				}
//				if cmd.Paths()[0] != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
//				}
//				if cmd.Paths()[1] != "/e/f/g" {
//					return fmt.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
//				}
//
//				return nil
//			},
//		},
//		{
//			name: "COPY no arg",
//			source: func() string {
//				return "COPY "
//			},
//			shouldFail: true,
//		},
//		{
//			name: "COPY with quoted default with ebedded colon",
//			source: func() string {
//				return `COPY '/a/:b/c'`
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 1 {
//					return fmt.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
//				}
//
//				arg := cmd.Paths()[0]
//				if arg != "/a/:b/c" {
//					return fmt.Errorf("COPY has unexpected argument %s", arg)
//				}
//				return nil
//			},
//		},
//		{
//			name: "COPY multiple with named param",
//			source: func() string {
//				return `COPY paths:"/a/b/c /e/:f/g"`
//			},
//			script: func(s *Script) error {
//				if len(s.Actions) != 1 {
//					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
//				}
//
//				cmd := s.Actions[0].(*CopyCommand)
//				if len(cmd.Paths()) != 2 {
//					return fmt.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
//				}
//				if cmd.Paths()[0] != "/a/b/c" {
//					return fmt.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
//				}
//				if cmd.Paths()[1] != "/e/:f/g" {
//					return fmt.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
//				}
//
//				return nil
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			runCommandTest(t, test)
//		})
//	}
//}

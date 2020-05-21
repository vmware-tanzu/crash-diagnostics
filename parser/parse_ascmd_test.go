// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package parser

//
//import (
//	"fmt"
//	"os"
//	"testing"
//
//	"github.com/vmware-tanzu/crash-diagnostics/script"
//)
//
//func TestCommandAS(t *testing.T) {
//	tests := []commandTest{
//		{
//			name: "AS specified with userid and groupid",
//			source: func() string {
//				return "AS userid:foo groupid:bar"
//			},
//			script: func(s *script.Script) error {
//				cmds := s.Preambles[CmdAs]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing preamble %s", CmdAs)
//				}
//				asCmd, ok := cmds[0].(*AsCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if asCmd.GetUserId() != "foo" {
//					return fmt.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
//				}
//				if asCmd.GetGroupId() != "bar" {
//					return fmt.Errorf("Unexpected AS groupid %s", asCmd.GetUserId())
//				}
//				return nil
//			},
//		},
//		{
//			name: "AS with quoted userid and groupid",
//			source: func() string {
//				return `AS userid:"foo" groupid:bar`
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAs]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing preamble %s", CmdAs)
//				}
//				asCmd, ok := cmds[0].(*AsCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if asCmd.GetUserId() != "foo" {
//					return fmt.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
//				}
//				if asCmd.GetGroupId() != "bar" {
//					return fmt.Errorf("Unexpected AS groupid %s", asCmd.GetUserId())
//				}
//				return nil
//			},
//		},
//		{
//			name: "AS with only userid",
//			source: func() string {
//				return "AS userid:foo"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAs]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing preamble %s", CmdAs)
//				}
//				asCmd, ok := cmds[0].(*AsCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if asCmd.GetUserId() != "foo" {
//					return fmt.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
//				}
//				if asCmd.GetGroupId() != fmt.Sprintf("%d", os.Getgid()) {
//					return fmt.Errorf("Unexpected AS groupid %s", asCmd.GetGroupId())
//				}
//				return nil
//			},
//		},
//		{
//			name: "AS not specified",
//			source: func() string {
//				return "FROM local"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAs]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing default AS preamble")
//				}
//				asCmd, ok := cmds[0].(*AsCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if asCmd.GetUserId() != fmt.Sprintf("%d", os.Getuid()) {
//					return fmt.Errorf("Unexpected AS default userid %s", asCmd.GetUserId())
//				}
//				if asCmd.GetGroupId() != fmt.Sprintf("%d", os.Getgid()) {
//					return fmt.Errorf("Unexpected AS default groupid %s", asCmd.GetUserId())
//				}
//				return nil
//			},
//		},
//		{
//			name: "Multiple AS provided",
//			source: func() string {
//				return "AS userid:foo\nAS userid:bar"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAs]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script should only have 1 AS instruction, got %d", len(cmds))
//				}
//				asCmd := cmds[0].(*AsCommand)
//				if asCmd.GetUserId() != "bar" {
//					return fmt.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
//				}
//				if asCmd.GetGroupId() != "" {
//					return fmt.Errorf("Unexpected AS groupid %s", asCmd.GetUserId())
//				}
//				return nil
//			},
//			shouldFail: true,
//		},
//		{
//			name: "AS with var expansion",
//			source: func() string {
//				os.Setenv("foogid", "barid")
//				return "AS userid:$USER groupid:$foogid"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAs]
//				asCmd := cmds[0].(*AsCommand)
//				if asCmd.GetUserId() != ExpandEnv("$USER") {
//					return fmt.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
//				}
//				if asCmd.GetGroupId() != "barid" {
//					return fmt.Errorf("Unexpected AS groupid %s", asCmd.GetUserId())
//				}
//				return nil
//			},
//		},
//		{
//			name: "AS with multiple args",
//			source: func() string {
//				return "AS foo:bar fuzz:buzz"
//			},
//			shouldFail: true,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			runCommandTest(t, test)
//		})
//	}
//}

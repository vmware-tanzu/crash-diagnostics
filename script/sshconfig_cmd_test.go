// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"testing"
)

func TestCommandSSHCONFIG(t *testing.T) {
	tests := []commandTest{
		{
			name: "SSHCONFIG specified with userid and private key",
			source: func() string {
				return "SSHCONFIG foo:/a/b/c"
			},
			script: func(s *Script) error {
				cmds := s.Preambles[CmdSSHConfig]
				if len(cmds) != 1 {
					return fmt.Errorf("Script missing preamble %s", CmdSSHConfig)
				}
				sshCmd, ok := cmds[0].(*SSHConfigCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", cmds[0])
				}
				if sshCmd.GetUserId() != "foo" {
					return fmt.Errorf("Unexpected AS userid %s", sshCmd.GetUserId())
				}
				if sshCmd.GetPrivateKeyPath() != "/a/b/c" {
					return fmt.Errorf("Unexpected AS groupid %s", sshCmd.GetUserId())
				}
				return nil
			},
		},
		{
			name: "SSHCONFIG with only private keypath",
			source: func() string {
				return "SSHCONFIG /a/b/c"
			},
			script: func(s *Script) error {
				cmds := s.Preambles[CmdSSHConfig]
				if len(cmds) != 1 {
					return fmt.Errorf("Script missing preamble %s", CmdSSHConfig)
				}
				sshCmd, ok := cmds[0].(*SSHConfigCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", cmds[0])
				}
				if sshCmd.GetUserId() != "" {
					return fmt.Errorf("Unexpected AS userid %s", sshCmd.GetUserId())
				}
				if sshCmd.GetPrivateKeyPath() != "/a/b/c" {
					return fmt.Errorf("Unexpected AS groupid %s", sshCmd.GetUserId())
				}
				return nil
			},
		},
		{
			name: "Multiple SSHCONFIG provided",
			source: func() string {
				return "SSHCONFIG /foo/bar\nSSHCONFIG bar:/bar"
			},
			script: func(s *Script) error {
				return nil
			},
			shouldFail: true,
		},
		{
			name: "SSHCONFIG with multiple args",
			source: func() string {
				return "SSHCONFIG foo:bar buzz"
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

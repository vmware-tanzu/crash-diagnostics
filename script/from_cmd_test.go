// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"os"
	"testing"
)

func TestCommandFROM(t *testing.T) {
	tests := []commandTest{
		{
			name: "FROM set to local",
			source: func() string {
				return "FROM local"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd, ok := froms[0].(*FromCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", froms[0])
				}
				if len(fromCmd.Nodes()) != 1 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Nodes()))
				}
				m := fromCmd.Nodes()[0]
				if m.Address() != Defaults.LocalSSHAddr {
					return fmt.Errorf("FROM has unexpected host address %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM set to single remote machine",
			source: func() string {
				return "FROM 'foo.bar:1234'"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)

				if len(fromCmd.Nodes()) != 1 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Nodes()))
				}
				m := fromCmd.Nodes()[0]
				if m.Address() != "foo.bar:1234" {
					return fmt.Errorf("FROM has unexpected machine address %s", m.Address())
				}
				if m.Address() != "foo.bar:1234" {
					return fmt.Errorf("FROM has unexpected machine host value %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM with multiple machines",
			source: func() string {
				return "FROM 'local 127.0.0.1:1234'"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)
				if len(fromCmd.Nodes()) != 2 {
					t.Log("Nodes:", fromCmd.Nodes())
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Nodes()))
				}
				m0 := fromCmd.Nodes()[0]
				m1 := fromCmd.Nodes()[1]
				if m0.Address() != Defaults.LocalSSHAddr {
					return fmt.Errorf("FROM arg 0 has unexpected host value: %s", m0.Address())
				}
				if m1.Address() != "127.0.0.1:1234" {
					return fmt.Errorf("FROM arg 1 has unexpected address:%s", m1.Address())
				}

				return nil
			},
		},
		{
			name: "Multiple FROMs",
			source: func() string {
				return "FROM local\nFROM local.1:123 local.2:456"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)
				if len(fromCmd.Nodes()) != 2 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Nodes()))
				}
				m0 := fromCmd.Nodes()[0]
				if m0.Address() != "local.1:123" {
					return fmt.Errorf("FROM arg 0 has unexpected machine arguments: %s", m0.Address())
				}
				m1 := fromCmd.Nodes()[1]
				if m1.Address() != "local.2:456" {
					return fmt.Errorf("FROM arg 1 has unexpected machine arguments: %s", m1.Address())
				}
				return nil
			},
		},
		{
			name: "FROM with named param",
			source: func() string {
				return "FROM hosts:local"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd, ok := froms[0].(*FromCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", froms[0])
				}
				if len(fromCmd.Nodes()) != 1 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Nodes()))
				}
				m := fromCmd.Nodes()[0]
				if m.Address() != Defaults.LocalSSHAddr {
					return fmt.Errorf("FROM has unexpected machine %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM remote machine named param",
			source: func() string {
				return "FROM hosts:foo.bar:1234"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)

				if len(fromCmd.Nodes()) != 1 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Nodes()))
				}
				m := fromCmd.Nodes()[0]
				if m.Address() != "foo.bar:1234" {
					return fmt.Errorf("FROM has unexpected machine address %s", m.Address())
				}
				if m.Address() != "foo.bar:1234" {
					return fmt.Errorf("FROM has unexpected machine host value %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM with var expansion",
			source: func() string {
				os.Setenv("foohost", "foo.bar")
				os.Setenv("fooport", "1234")
				return "FROM ${foohost}:$fooport"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)

				if len(fromCmd.Nodes()) != 1 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Nodes()))
				}
				m := fromCmd.Nodes()[0]
				if m.Address() != "foo.bar:1234" {
					return fmt.Errorf("FROM has unexpected machine address %s", m.Address())
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}

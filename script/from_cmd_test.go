package script

import (
	"fmt"
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
				if len(fromCmd.Machines()) != 1 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Machines()))
				}
				m := fromCmd.Machines()[0]
				if m.Host() != "local" {
					return fmt.Errorf("FROM has unexpected machine %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM set to single remote machine",
			source: func() string {
				return "FROM foo.bar:1234"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)

				if len(fromCmd.Machines()) != 1 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Machines()))
				}
				m := fromCmd.Machines()[0]
				if m.Address() != "foo.bar:1234" {
					return fmt.Errorf("FROM has unexpected machine address %s", m.Address())
				}
				if m.Host() != "foo.bar" {
					return fmt.Errorf("FROM has unexpected machine host value %s", m)
				}
				if m.Port() != "1234" {
					return fmt.Errorf("FROM has unexpected machine port value %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM with multiple machines",
			source: func() string {
				return "FROM local 127.0.0.1:1234"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)
				if len(fromCmd.Machines()) != 2 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Machines()))
				}
				m0 := fromCmd.Machines()[0]
				m1 := fromCmd.Machines()[1]
				if m0.Host() != "local" {
					return fmt.Errorf("FROM arg 0 has unexpected host value: %s", m0.Host())
				}
				if m1.Host() != "127.0.0.1" || m1.Port() != "1234" {
					return fmt.Errorf("FROM arg 1 has unexpected host:%s port:%s values", m1.Host(), m1.Port())
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
				if len(fromCmd.Machines()) != 2 {
					return fmt.Errorf("FROM has unexpected number of machines %d", len(fromCmd.Machines()))
				}
				m0 := fromCmd.Machines()[0]
				if m0.Host() != "local.1" || m0.Port() != "123" {
					return fmt.Errorf("FROM arg 0 has unexpected machine arguments: %s:%s", m0.Host(), m0.Port())
				}
				m1 := fromCmd.Machines()[1]
				if m1.Host() != "local.2" || m1.Port() != "456" {
					return fmt.Errorf("FROM arg 1 has unexpected machine arguments: %s:%s", m1.Host(), m1.Port())
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

package script

import (
	"fmt"
	"testing"
)

func TestCommandFROM(t *testing.T) {
	tests := []commandTest{
		{
			name: "FROM with single arg",
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
				if m.Address() != "local" {
					return fmt.Errorf("FROM has unexpected machine %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM with multiple machines",
			source: func() string {
				return "FROM local 127.0.0.1"
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
				if fromCmd.Machines()[0].Address() != "local" || fromCmd.Machines()[1].Address() != "127.0.0.1" {
					return fmt.Errorf("FROM has unexpected machine arguments: %s", fromCmd.Args())
				}
				return nil
			},
		},
		{
			name: "Multiple FROMs",
			source: func() string {
				return "FROM local\nFROM local.1 local.2"
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
				if fromCmd.Machines()[0].Address() != "local.1" || fromCmd.Machines()[1].Address() != "local.2" {
					return fmt.Errorf("FROM has unexpected machine arguments: %s", fromCmd.Args())
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

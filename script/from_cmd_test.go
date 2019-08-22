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
				if m.Address != "local" {
					return fmt.Errorf("FROM has unexpected machine %s", m)
				}
				return nil
			},
		},
		{
			name: "FROM with multiple machines",
			source: func() string {
				return "FROM local local"
			},
			shouldFail: true,
		},
		{
			name: "Multiple FROMs",
			source: func() string {
				return "FROM local\nFROM local2"
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

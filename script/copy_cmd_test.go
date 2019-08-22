package script

import (
	"fmt"
	"testing"
)

func TestCommandCOPY(t *testing.T) {
	tests := []commandTest{
		{
			name: "COPY with single arg",
			source: func() string {
				return "COPY /a/b/c"
			},
			script: func(s *Script) error {
				if len(s.Actions) != 1 {
					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
				}

				if len(s.Actions[0].Args()) != 1 {
					return fmt.Errorf("COPY has unexpected number of args %d", len(s.Actions[0].Args()))
				}

				arg := s.Actions[0].Args()[0]
				if arg != "/a/b/c" {
					return fmt.Errorf("COPY has unexpected argument %s", arg)
				}
				return nil
			},
		},
		{
			name: "COPY with multiple args",
			source: func() string {
				return "COPY /a/b/c /e/f/g"
			},
			script: func(s *Script) error {
				if len(s.Actions) != 1 {
					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
				}

				cmd := s.Actions[0]
				if len(cmd.Args()) != 2 {
					return fmt.Errorf("COPY has unexpected number of args %d", len(cmd.Args()))
				}
				if cmd.Args()[0] != "/a/b/c" {
					return fmt.Errorf("COPY has unexpected argument[0] %s", cmd.Args()[0])
				}
				if cmd.Args()[1] != "/e/f/g" {
					return fmt.Errorf("COPY has unexpected argument[1] %s", cmd.Args()[1])
				}

				return nil
			},
		},
		{
			name: "Multiple COPY commands",
			source: func() string {
				return "COPY /a/b/c\nCOPY d /e/f"
			},
			script: func(s *Script) error {
				if len(s.Actions) != 2 {
					return fmt.Errorf("Script has unexpected COPY actions, has %d COPY", len(s.Actions))
				}

				if len(s.Actions[0].Args()) != 1 {
					return fmt.Errorf("COPY action[0] has wrong number of args %d", len(s.Actions[0].Args()))
				}
				if len(s.Actions[1].Args()) != 2 {
					return fmt.Errorf("COPY action[1] has wrong number of args %d", len(s.Actions[1].Args()))
				}
				arg := s.Actions[0].Args()[0]
				if arg != "/a/b/c" {
					return fmt.Errorf("COPY action[0] has unexpected arg %s", arg)
				}
				arg = s.Actions[1].Args()[0]
				if arg != "d" {
					return fmt.Errorf("COPY action[1] has unexpected arg[0] %s", arg)
				}
				arg = s.Actions[1].Args()[1]
				if arg != "/e/f" {
					return fmt.Errorf("COPY action[1] has unexpected arg[1] %s", arg)
				}
				return nil
			},
		},
		{
			name: "COPY no arg",
			source: func() string {
				return "COPY "
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

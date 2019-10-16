// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"testing"
)

func TestCommandCAPTURE(t *testing.T) {
	tests := []commandTest{
		{
			name: "CAPTURE single command",
			source: func() string {
				return "CAPTURE /bin/echo HELLO WORLD"
			},
			script: func(s *Script) error {
				if len(s.Actions) != 1 {
					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
				}
				cmd, ok := s.Actions[0].(*CaptureCommand)
				if !ok {
					return fmt.Errorf("Unexpected action type %T in script", s.Actions[0])
				}

				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}

				cliCmd, cliArgs := cmd.GetParsedCmd()
				if cliCmd != "/bin/echo" {
					return fmt.Errorf("CAPTURE action parsed cli unexpected command %s", cliCmd)
				}
				if len(cliArgs) != 2 {
					return fmt.Errorf("CAPTURE action parsed cli unexpected args %d", len(cliArgs))
				}

				return nil
			},
		},
		{
			name: "CAPTURE multiple commands",
			source: func() string {
				return "CAPTURE /bin/echo HELLO\nCOPY a/b\nCAPTURE /bin/clear"
			},
			script: func(s *Script) error {
				if len(s.Actions) != 3 {
					return fmt.Errorf("Script has unexpected number of actions: %d", len(s.Actions))
				}
				cmd0, ok := s.Actions[0].(*CaptureCommand)
				if !ok {
					return fmt.Errorf("Unexpected action type %T at pos 0", s.Actions[0])
				}
				cmd2, ok := s.Actions[2].(*CaptureCommand)
				if !ok {
					return fmt.Errorf("Unexpected action type %T at pos 0", s.Actions[2])
				}

				if cmd0.Args()["cmd"] != cmd0.GetCmdString() {
					return fmt.Errorf("CAPTURE action 0 with unexpected CLI string %s", cmd0.GetCmdString())
				}
				if cmd2.Args()["cmd"] != cmd2.GetCmdString() {
					return fmt.Errorf("CAPTURE action 2 with unexpected CLI string %s", cmd2.GetCmdString())
				}
				return nil
			},
		},
		{
			name: "CAPTURE with named parameters",
			source: func() string {
				return "CAPTURE cmd:/bin/date"
			},
			script: func(s *Script) error {
				if len(s.Actions) != 1 {
					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
				}
				cmd, ok := s.Actions[0].(*CaptureCommand)
				if !ok {
					return fmt.Errorf("Unexpected action type %T in script", s.Actions[0])
				}

				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}

				cliCmd, cliArgs := cmd.GetParsedCmd()
				if cliCmd != "/bin/date" {
					return fmt.Errorf("CAPTURE action parsed cli unexpected command %s", cliCmd)
				}
				if len(cliArgs) != 0 {
					return fmt.Errorf("CAPTURE action parsed cli unexpected args %d", len(cliArgs))
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

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"os"
	"testing"
)

func TestCommandRUN(t *testing.T) {
	tests := []commandTest{
		{
			name: "RUN with default param",
			source: func() string {
				return `RUN /bin/echo "HELLO WORLD"`
			},
			script: func(s *Script) error {
				if len(s.Actions) != 1 {
					return fmt.Errorf("Script has unexpected action count, needs %d", len(s.Actions))
				}
				cmd, ok := s.Actions[0].(*RunCommand)
				if !ok {
					return fmt.Errorf("Unexpected action type %T in script", s.Actions[0])
				}

				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					return fmt.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					return fmt.Errorf("RUN command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					return fmt.Errorf("RUN unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 1 {
					return fmt.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
				}

				return nil
			},
		},
		{
			name: "RUN default quoted param",
			source: func() string {
				return `RUN '/bin/echo "HELLO WORLD"'`
			},
			script: func(s *Script) error {
				if len(s.Actions) != 1 {
					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
				}
				cmd := s.Actions[0].(*RunCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					return fmt.Errorf("RUN action with unexpected CLI string %s", cmd.GetCmdString())
				}

				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					return fmt.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					return fmt.Errorf("RUN command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					return fmt.Errorf("RUN unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 1 {
					return fmt.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
				}

				return nil
			},
		},
		{
			name: "RUN named param command",
			source: func() string {
				return `RUN cmd:"/bin/echo 'HELLO WORLD'"`
			},
			script: func(s *Script) error {
				cmd := s.Actions[0].(*RunCommand)

				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					return fmt.Errorf("CAPTURE command parse failed: %s", cmd.GetCmdString())
				}
				if cliCmd != "/bin/echo" {
					return fmt.Errorf("CAPTURE action parsed cli unexpected command %s", cliCmd)
				}
				if len(cliArgs) != 1 {
					return fmt.Errorf("CAPTURE action parsed cli unexpected args %d", len(cliArgs))
				}

				return nil
			},
		},
		{
			name: "RUN multiple commands",
			source: func() string {
				return `
				RUN /bin/echo "HELLO WORLD"
				COPY a/b
				RUN cmd:"/bin/clear"`
			},
			script: func(s *Script) error {
				if len(s.Actions) != 3 {
					return fmt.Errorf("Script has unexpected number of actions: %d", len(s.Actions))
				}
				cmd0 := s.Actions[0].(*RunCommand)
				cmd2 := s.Actions[2].(*RunCommand)
				if cmd0.Args()["cmd"] != cmd0.GetCmdString() {
					return fmt.Errorf("RUN at 0 with unexpected command string %s", cmd0.GetCmdString())
				}
				if cmd2.Args()["cmd"] != cmd2.GetCmdString() {
					return fmt.Errorf("RUN at 2 with unexpected command string %s", cmd2.GetCmdString())
				}
				return nil
			},
		},
		{
			name: "RUN with named parameters",
			source: func() string {
				return "RUN cmd:/bin/date"
			},
			script: func(s *Script) error {
				cmd := s.Actions[0].(*RunCommand)
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					return fmt.Errorf("RUN command parse failed: %s", err)
				}
				if cliCmd != "/bin/date" {
					return fmt.Errorf("RUN parsed unexpected command name: %s", cliCmd)
				}
				if len(cliArgs) != 0 {
					return fmt.Errorf("RUN parsed unexpected command args: %d", len(cliArgs))
				}

				return nil
			},
		},
		{
			name: "RUN with expanded vars",
			source: func() string {
				os.Setenv("msg", "Hello World!")
				return `RUN '/bin/echo "$msg"'`
			},
			script: func(s *Script) error {
				cmd := s.Actions[0].(*RunCommand)
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					return fmt.Errorf("RUN command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					return fmt.Errorf("RUN parsed unexpected command name %s", cliCmd)
				}
				if cliArgs[0] != "Hello World!" {
					return fmt.Errorf("RUN parsed unexpected command args: %s", cliArgs)
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

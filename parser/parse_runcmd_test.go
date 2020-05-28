// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"testing"
)

func TestCommandRUN(t *testing.T) {
	tests := []parserTest{
		//{
		//	name: "RUN",
		//	source: func(t *testing.T) string {
		//		return `RUN /bin/echo "HELLO WORLD"`
		//	},
		//	script: func(t *testing.T, s *script.Script) {
		//		if len(s.Actions) != 1 {
		//			t.Errorf("Script has unexpected action count, needs %d", len(s.Actions))
		//		}
		//		cmd, ok := s.Actions[0].(*script.RunCommand)
		//		if !ok {
		//			t.Fatalf("Unexpected action type %T in script", s.Actions[0])
		//		}
		//
		//		if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//			t.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
		//		}
		//		cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//		if err != nil {
		//			t.Errorf("RUN command parse failed: %s", err)
		//		}
		//		if cliCmd != "/bin/echo" {
		//			t.Errorf("RUN unexpected command parsed: %s", cliCmd)
		//		}
		//		if len(cliArgs) != 1 {
		//			t.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
		//		}
		//		if cliArgs[0] != "HELLO WORLD" {
		//			t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//		}
		//	},
		//},
		//		{
		//			name: "RUN single-quoted default with quoted param",
		//			source: func() string {
		//				return `RUN '/bin/echo -n "HELLO WORLD"'`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*RunCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("RUN action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("RUN unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO WORLD" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN single-quoted named param with quoted arg",
		//			source: func() string {
		//				return `RUN cmd:'/bin/echo -n "HELLO WORLD"'`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*RunCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("RUN action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("RUN unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO WORLD" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN double-quoted named param with quoted arg",
		//			source: func() string {
		//				return `RUN cmd:"/bin/echo -n 'HELLO WORLD'"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*RunCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("RUN action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("RUN unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO WORLD" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN multiple commands",
		//			source: func() string {
		//				return `
		//				RUN /bin/echo "HELLO WORLD"
		//				RUN cmd:"/bin/bash -c date"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 2 {
		//					return fmt.Errorf("Script has unexpected number of actions: %d", len(s.Actions))
		//				}
		//				cmd0 := s.Actions[0].(*RunCommand)
		//				cmd2 := s.Actions[1].(*RunCommand)
		//				if cmd0.Args()["cmd"] != cmd0.GetCmdString() {
		//					return fmt.Errorf("RUN at 0 with unexpected command string %s", cmd0.GetCmdString())
		//				}
		//				if cmd2.Args()["cmd"] != cmd2.GetCmdString() {
		//					return fmt.Errorf("RUN at 2 with unexpected command string %s", cmd2.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd2.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/bash" {
		//					return fmt.Errorf("RUN unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-c" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "date" {
		//					return fmt.Errorf("RUN has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN unquoted named params",
		//			source: func() string {
		//				return "RUN cmd:/bin/date"
		//			},
		//			script: func(s *Script) error {
		//				cmd := s.Actions[0].(*RunCommand)
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/date" {
		//					return fmt.Errorf("RUN parsed unexpected command name: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 0 {
		//					return fmt.Errorf("RUN parsed unexpected command args: %d", len(cliArgs))
		//				}
		//
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN with expanded vars",
		//			source: func() string {
		//				os.Setenv("msg", "Hello World!")
		//				return `RUN '/bin/echo "$msg"'`
		//			},
		//			script: func(s *Script) error {
		//				cmd := s.Actions[0].(*RunCommand)
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("RUN parsed unexpected command name %s", cliCmd)
		//				}
		//				if cliArgs[0] != "Hello World!" {
		//					return fmt.Errorf("RUN parsed unexpected command args: %s", cliArgs)
		//				}
		//
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN unquoted default with quoted subproc",
		//			source: func() string {
		//				return `RUN /bin/bash -c 'echo "Hello World"'`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*RunCommand)
		//				effCmd, err := cmd.GetEffectiveCmdStr()
		//				if err != nil {
		//					return fmt.Errorf("RUN effective command str failed: %s", err)
		//				}
		//				if effCmd != `/bin/bash -c 'echo "Hello World"'` {
		//					return fmt.Errorf("RUN unexpected effective command str: %s", effCmd)
		//				}
		//
		//				effArgs, err := cmd.GetEffectiveCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN effective command args failed: %s", err)
		//				}
		//				if len(effArgs) != 3 {
		//					return fmt.Errorf("RUN unexpected effective command args: %#v", effArgs)
		//				}
		//
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("RUN unexpected command args parsed: %#v", cliArgs)
		//				}
		//				if cliCmd != "/bin/bash" {
		//					return fmt.Errorf("RUN unexpected command parsed: %#v", cliCmd)
		//				}
		//				if strings.TrimSpace(cliArgs[0]) != "-c" {
		//					return fmt.Errorf("RUN has unexpected shell argument: expecting -c, got %#v", cliArgs)
		//				}
		//				if cliArgs[1] != `echo "Hello World"` {
		//					return fmt.Errorf("RUN has unexpected subproc argument: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN quoted with shell and quoted subproc",
		//			source: func() string {
		//				return `RUN shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*RunCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
		//				}
		//				if cmd.Args()["shell"] != cmd.GetCmdShell() {
		//					return fmt.Errorf("RUN action with unexpected shell %s", cmd.GetCmdShell())
		//				}
		//				effCmdStr, err := cmd.GetEffectiveCmdStr()
		//				if err != nil {
		//					return fmt.Errorf("RUN effective command str failed: %s", err)
		//				}
		//				if effCmdStr != `/bin/bash -c "echo 'HELLO WORLD'"` {
		//					return fmt.Errorf("RUN unexpected effective command string: %s", effCmdStr)
		//				}
		//
		//				effArgs, err := cmd.GetEffectiveCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN effective command args failed: %s", err)
		//				}
		//				if len(effArgs) != 3 {
		//					return fmt.Errorf("RUN unexpected effective command args: %#v", effArgs)
		//				}
		//
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("RUN command parse failed: %s", err)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("RUN unexpected command args parsed: %#v", cliArgs)
		//				}
		//				if cliCmd != "/bin/bash" {
		//					return fmt.Errorf("RUN unexpected command parsed: %#v", cliCmd)
		//				}
		//				if cliArgs[0] != "-c" {
		//					return fmt.Errorf("RUN has unexpected shell argument: expecting -c, got %s", cliArgs[0])
		//				}
		//				if cliArgs[1] != "echo 'HELLO WORLD'" {
		//					return fmt.Errorf("RUN has unexpected shell argument: %s", cliArgs[0])
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "RUN with echo param",
		//			source: func() string {
		//				return `RUN shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'" echo:"true"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*RunCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
		//				}
		//				if cmd.Args()["shell"] != cmd.GetCmdShell() {
		//					return fmt.Errorf("RUN action with unexpected shell %s", cmd.GetCmdShell())
		//				}
		//				if cmd.Args()["echo"] != cmd.GetEcho() {
		//					return fmt.Errorf("RUN action with unexpected echo param %s", cmd.GetCmdShell())
		//				}
		//				return nil
		//			},
		//		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runParserTest(t, test)
		})
	}
}

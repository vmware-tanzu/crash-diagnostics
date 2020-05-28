// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"testing"
)

func TestCommandRUN(t *testing.T) {
	//tests := []commandTest{
	//	{
	//		name: "RUN",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, `/bin/echo "HELLO WORLD"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd, ok := c.(*RunCommand)
	//			if !ok {
	//				t.Errorf("Unexpected action type %T in script", c)
	//			}
	//
	//			if cmd.Args()["cmd"] != cmd.GetCmdString() {
	//				t.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
	//			}
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if cliCmd != "/bin/echo" {
	//				t.Errorf("RUN unexpected command parsed: %s", cliCmd)
	//			}
	//			if len(cliArgs) != 1 {
	//				t.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
	//			}
	//			if cliArgs[0] != "HELLO WORLD" {
	//				t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/single quoted",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, `'/bin/echo -n "HELLO WORLD"'`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			if cmd.Args()["cmd"] != cmd.GetCmdString() {
	//				t.Errorf("RUN action with unexpected CLI string %s", cmd.GetCmdString())
	//			}
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if cliCmd != "/bin/echo" {
	//				t.Errorf("RUN unexpected command parsed: %s", cliCmd)
	//			}
	//			if len(cliArgs) != 2 {
	//				t.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
	//			}
	//			if cliArgs[0] != "-n" {
	//				t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
	//			}
	//			if cliArgs[1] != "HELLO WORLD" {
	//				t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/param single quoted",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, `cmd:'/bin/echo -n "HELLO WORLD"'`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			if cmd.Args()["cmd"] != cmd.GetCmdString() {
	//				t.Errorf("RUN action with unexpected CLI string %s", cmd.GetCmdString())
	//			}
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if cliCmd != "/bin/echo" {
	//				t.Errorf("RUN unexpected command parsed: %s", cliCmd)
	//			}
	//			if len(cliArgs) != 2 {
	//				t.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
	//			}
	//			if cliArgs[0] != "-n" {
	//				t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
	//			}
	//			if cliArgs[1] != "HELLO WORLD" {
	//				t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/cmd double-quoted",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, `cmd:"/bin/echo -n 'HELLO WORLD'"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			if cmd.Args()["cmd"] != cmd.GetCmdString() {
	//				t.Errorf("RUN action with unexpected CLI string %s", cmd.GetCmdString())
	//			}
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if cliCmd != "/bin/echo" {
	//				t.Errorf("RUN unexpected command parsed: %s", cliCmd)
	//			}
	//			if len(cliArgs) != 2 {
	//				t.Errorf("RUN unexpected command args parsed: %d", len(cliArgs))
	//			}
	//			if cliArgs[0] != "-n" {
	//				t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
	//			}
	//			if cliArgs[1] != "HELLO WORLD" {
	//				t.Errorf("RUN has unexpected cli args: %#v", cliArgs)
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/cmd unquoted",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, "cmd:/bin/date")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if cliCmd != "/bin/date" {
	//				t.Errorf("RUN parsed unexpected command name: %s", cliCmd)
	//			}
	//			if len(cliArgs) != 0 {
	//				t.Errorf("RUN parsed unexpected command args: %d", len(cliArgs))
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/expanded vars",
	//		command: func(t *testing.T) Directive {
	//			os.Setenv("msg", "Hello World!")
	//			cmd, err := NewRunCommand(0, `'/bin/echo "$msg"'`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if cliCmd != "/bin/echo" {
	//				t.Errorf("RUN parsed unexpected command name %s", cliCmd)
	//			}
	//			if cliArgs[0] != "Hello World!" {
	//				t.Errorf("RUN parsed unexpected command args: %s", cliArgs)
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/multi quotes",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, `/bin/bash -c 'echo "Hello World"'`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			effCmd, err := cmd.GetEffectiveCmdStr()
	//			if err != nil {
	//				t.Errorf("RUN effective command str failed: %s", err)
	//			}
	//			if effCmd != `/bin/bash -c 'echo "Hello World"'` {
	//				t.Errorf("RUN unexpected effective command str: %s", effCmd)
	//			}
	//
	//			effArgs, err := cmd.GetEffectiveCmd()
	//			if err != nil {
	//				t.Errorf("RUN effective command args failed: %s", err)
	//			}
	//			if len(effArgs) != 3 {
	//				t.Errorf("RUN unexpected effective command args: %#v", effArgs)
	//			}
	//
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if len(cliArgs) != 2 {
	//				t.Errorf("RUN unexpected command args parsed: %#v", cliArgs)
	//			}
	//			if cliCmd != "/bin/bash" {
	//				t.Errorf("RUN unexpected command parsed: %#v", cliCmd)
	//			}
	//			if strings.TrimSpace(cliArgs[0]) != "-c" {
	//				t.Errorf("RUN has unexpected shell argument: expecting -c, got %#v", cliArgs)
	//			}
	//			if cliArgs[1] != `echo "Hello World"` {
	//				t.Errorf("RUN has unexpected subproc argument: %#v", cliArgs)
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/shell cmd quoted",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, `shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			if cmd.Args()["cmd"] != cmd.GetCmdString() {
	//				t.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
	//			}
	//			if cmd.Args()["shell"] != cmd.GetCmdShell() {
	//				t.Errorf("RUN action with unexpected shell %s", cmd.GetCmdShell())
	//			}
	//			effCmdStr, err := cmd.GetEffectiveCmdStr()
	//			if err != nil {
	//				t.Errorf("RUN effective command str failed: %s", err)
	//			}
	//			if effCmdStr != `/bin/bash -c "echo 'HELLO WORLD'"` {
	//				t.Errorf("RUN unexpected effective command string: %s", effCmdStr)
	//			}
	//
	//			effArgs, err := cmd.GetEffectiveCmd()
	//			if err != nil {
	//				t.Errorf("RUN effective command args failed: %s", err)
	//			}
	//			if len(effArgs) != 3 {
	//				t.Errorf("RUN unexpected effective command args: %#v", effArgs)
	//			}
	//
	//			cliCmd, cliArgs, err := cmd.GetParsedCmd()
	//			if err != nil {
	//				t.Errorf("RUN command parse failed: %s", err)
	//			}
	//			if len(cliArgs) != 2 {
	//				t.Errorf("RUN unexpected command args parsed: %#v", cliArgs)
	//			}
	//			if cliCmd != "/bin/bash" {
	//				t.Errorf("RUN unexpected command parsed: %#v", cliCmd)
	//			}
	//			if cliArgs[0] != "-c" {
	//				t.Errorf("RUN has unexpected shell argument: expecting -c, got %s", cliArgs[0])
	//			}
	//			if cliArgs[1] != "echo 'HELLO WORLD'" {
	//				t.Errorf("RUN has unexpected shell argument: %s", cliArgs[0])
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "RUN/echo",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewRunCommand(0, `shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'" echo:"true"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cmd := c.(*RunCommand)
	//			if cmd.Args()["cmd"] != cmd.GetCmdString() {
	//				t.Errorf("RUN action with unexpected command string %s", cmd.GetCmdString())
	//			}
	//			if cmd.Args()["shell"] != cmd.GetCmdShell() {
	//				t.Errorf("RUN action with unexpected shell %s", cmd.GetCmdShell())
	//			}
	//			if cmd.Args()["echo"] != cmd.GetEcho() {
	//				t.Errorf("RUN action with unexpected echo param %s", cmd.GetCmdShell())
	//			}
	//
	//		},
	//	},
	//}
	//
	//for _, test := range tests {
	//	t.Run(test.name, func(t *testing.T) {
	//		runCommandTest(t, test)
	//	})
	//}
}

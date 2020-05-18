// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"testing"
)

func TestCommandCAPTURE(t *testing.T) {
	tests := []commandTest{
		{
			name: "CAPTURE/single unquoted param",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `CAPTURE /bin/echo "HELLO WORLD"`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd, ok := c.(*CaptureCommand)
				if !ok {
					t.Fatalf("Unexpected action type %T in script", c)
				}

				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected command string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 1 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "HELLO WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/single quoted param",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `'/bin/echo -n "HELLO WORLD"'`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Fatalf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "-n" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
				if cliArgs[1] != "HELLO WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/single-quoted named param",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `cmd:'/bin/echo -n "HELLO WORLD"'`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "-n" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
				if cliArgs[1] != "HELLO WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/double-quoted named param",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `cmd:"/bin/echo -n 'HELLO WORLD'"`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "-n" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
				if cliArgs[1] != "HELLO WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/unquoted named params",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0,"cmd:/bin/date")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/date" {
					t.Errorf("CAPTURE parsed unexpected command name: %s", cliCmd)
				}
				if len(cliArgs) != 0 {
					t.Errorf("CAPTURE parsed unexpected command args: %d", len(cliArgs))
				}
			},
		},
		{
			name: "CAPTURE/expanded vars",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `'/bin/echo "$msg"'`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				os.Setenv("msg", "Hello World!")
				cmd := c.(*CaptureCommand)
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE parsed unexpected command name %s", cliCmd)
				}
				if cliArgs[0] != "Hello World!" {
					t.Errorf("CAPTURE parsed unexpected command args: %s", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/with shell",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'"`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected command string %s", cmd.GetCmdString())
				}
				if cmd.Args()["shell"] != cmd.GetCmdShell() {
					t.Errorf("CAPTURE action with unexpected shell %s", cmd.GetCmdShell())
				}

				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %#v", cliArgs)
				}
				if cliCmd != "/bin/bash" {
					t.Errorf("CAPTURE unexpected command parsed: %#v", cliCmd)
				}
				if cliArgs[0] != "-c" {
					t.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
				}
				if cliArgs[1] != "echo 'HELLO WORLD'" {
					t.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
				}
			},
		},
		{
			name: "CAPTURE/with echo",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0,`shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'" echo:"true"`)
				if err != nil{
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Fatalf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				if cmd.Args()["shell"] != cmd.GetCmdShell() {
					t.Errorf("CAPTURE action with unexpected shell %s", cmd.GetCmdShell())
				}
				if cmd.Args()["echo"] != cmd.GetEcho() {
					t.Errorf("CAPTURE action with unexpected echo param %s", cmd.GetCmdShell())
				}
			},
		},
		{
			name: "CAPTURE/unquote with embedded colons",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `/bin/echo "HELLO:WORLD"`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd, ok := c.(*CaptureCommand)
				if !ok {
					t.Errorf("Unexpected action type %T in script", c)
				}

				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected command string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 1 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "HELLO:WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/quoted with embedded colon",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `'/bin/echo -n "HELLO:WORLD"'`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "-n" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
				if cliArgs[1] != "HELLO:WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/single-quoted named-param with embedded colon",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `cmd:'/bin/echo -n "HELLO:WORLD"'`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "-n" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
				if cliArgs[1] != "HELLO:WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE/double-quoted named-param with embedded colon",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0,`cmd:"/bin/echo -n 'HELLO:WORLD'"`)
				if err != nil{
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
				}
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/echo" {
					t.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
				}
				if cliArgs[0] != "-n" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
				if cliArgs[1] != "HELLO:WORLD" {
					t.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
				}
			},
		},
		{
			name: "CAPTURE unquoted named param with multiple embedded colons",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, "cmd:/bin/date:time:")
				if err != nil{
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if cliCmd != "/bin/date:time:" {
					t.Errorf("CAPTURE parsed unexpected command name: %s", cliCmd)
				}
				if len(cliArgs) != 0 {
					t.Errorf("CAPTURE parsed unexpected command args: %d", len(cliArgs))
				}
			},
		},
		{
			name: "CAPTURE/shell with embedded colon",
			command: func(t *testing.T) Command {
				cmd, err := NewCaptureCommand(0, `shell:"/bin/bash -c" cmd:"echo 'HELLO:WORLD'"`)
				if err != nil{
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CaptureCommand)
				if cmd.Args()["cmd"] != cmd.GetCmdString() {
					t.Errorf("CAPTURE action with unexpected command string %s", cmd.GetCmdString())
				}
				if cmd.Args()["shell"] != cmd.GetCmdShell() {
					t.Errorf("CAPTURE action with unexpected shell %s", cmd.GetCmdShell())
				}

				cliCmd, cliArgs, err := cmd.GetParsedCmd()
				if err != nil {
					t.Errorf("CAPTURE command parse failed: %s", err)
				}
				if len(cliArgs) != 2 {
					t.Errorf("CAPTURE unexpected command args parsed: %#v", cliArgs)
				}
				if cliCmd != "/bin/bash" {
					t.Errorf("CAPTURE unexpected command parsed: %#v", cliCmd)
				}
				if cliArgs[0] != "-c" {
					t.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
				}
				if cliArgs[1] != "echo 'HELLO:WORLD'" {
					t.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}

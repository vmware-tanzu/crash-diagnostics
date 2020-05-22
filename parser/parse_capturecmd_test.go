// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestCommandCAPTURE(t *testing.T) {
	tests := []parserTest{
		{
			name: "CAPTURE",
			source: func(t *testing.T) string {
				return `CAPTURE /bin/echo "HELLO WORLD"`
			},
			script: func(t *testing.T, s *script.Script) {
				if len(s.Actions) != 1 {
					t.Errorf("Script has unexpected action count, needs %d", len(s.Actions))
				}
				cmd, ok := s.Actions[0].(*script.CaptureCommand)
				if !ok {
					t.Errorf("Unexpected action type %T in script", s.Actions[0])
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
		//		{
		//			name: "CAPTURE single-quoted default with quoted param",
		//			source: func() string {
		//				return `CAPTURE '/bin/echo -n "HELLO WORLD"'`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO WORLD" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE single-quoted named param with quoted arg",
		//			source: func() string {
		//				return `CAPTURE cmd:'/bin/echo -n "HELLO WORLD"'`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO WORLD" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE double-quoted named param with quoted arg",
		//			source: func() string {
		//				return `CAPTURE cmd:"/bin/echo -n 'HELLO WORLD'"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO WORLD" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE multiple commands",
		//			source: func() string {
		//				return `
		//				CAPTURE /bin/echo "HELLO WORLD"
		//				CAPTURE cmd:"/bin/bash -c date"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 2 {
		//					return fmt.Errorf("Script has unexpected number of actions: %d", len(s.Actions))
		//				}
		//				cmd0 := s.Actions[0].(*CaptureCommand)
		//				cmd2 := s.Actions[1].(*CaptureCommand)
		//				if cmd0.Args()["cmd"] != cmd0.GetCmdString() {
		//					return fmt.Errorf("CAPTURE at 0 with unexpected command string %s", cmd0.GetCmdString())
		//				}
		//				if cmd2.Args()["cmd"] != cmd2.GetCmdString() {
		//					return fmt.Errorf("CAPTURE at 2 with unexpected command string %s", cmd2.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd2.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/bash" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-c" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "date" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE unquoted named params",
		//			source: func() string {
		//				return "CAPTURE cmd:/bin/date"
		//			},
		//			script: func(s *Script) error {
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/date" {
		//					return fmt.Errorf("CAPTURE parsed unexpected command name: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 0 {
		//					return fmt.Errorf("CAPTURE parsed unexpected command args: %d", len(cliArgs))
		//				}
		//
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE with expanded vars",
		//			source: func() string {
		//				os.Setenv("msg", "Hello World!")
		//				return `CAPTURE '/bin/echo "$msg"'`
		//			},
		//			script: func(s *Script) error {
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE parsed unexpected command name %s", cliCmd)
		//				}
		//				if cliArgs[0] != "Hello World!" {
		//					return fmt.Errorf("CAPTURE parsed unexpected command args: %s", cliArgs)
		//				}
		//
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE quoted with shell and quoted subproc",
		//			source: func() string {
		//				return `CAPTURE shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected command string %s", cmd.GetCmdString())
		//				}
		//				if cmd.Args()["shell"] != cmd.GetCmdShell() {
		//					return fmt.Errorf("CAPTURE action with unexpected shell %s", cmd.GetCmdShell())
		//				}
		//
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %#v", cliArgs)
		//				}
		//				if cliCmd != "/bin/bash" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %#v", cliCmd)
		//				}
		//				if cliArgs[0] != "-c" {
		//					return fmt.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
		//				}
		//				if cliArgs[1] != "echo 'HELLO WORLD'" {
		//					return fmt.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE with echo param",
		//			source: func() string {
		//				return `CAPTURE shell:"/bin/bash -c" cmd:"echo 'HELLO WORLD'" echo:"true"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				if cmd.Args()["shell"] != cmd.GetCmdShell() {
		//					return fmt.Errorf("CAPTURE action with unexpected shell %s", cmd.GetCmdShell())
		//				}
		//				if cmd.Args()["echo"] != cmd.GetEcho() {
		//					return fmt.Errorf("CAPTURE action with unexpected echo param %s", cmd.GetCmdShell())
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE with unqoted default with embeded colons",
		//			source: func() string {
		//				return `CAPTURE /bin/echo "HELLO:WORLD"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected action count, needs %d", len(s.Actions))
		//				}
		//				cmd, ok := s.Actions[0].(*CaptureCommand)
		//				if !ok {
		//					return fmt.Errorf("Unexpected action type %T in script", s.Actions[0])
		//				}
		//
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected command string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 1 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "HELLO:WORLD" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE single-quoted-default with embedded colon",
		//			source: func() string {
		//				return `CAPTURE '/bin/echo -n "HELLO:WORLD"'`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO:WORLD" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE single-quoted named param with embedded colon",
		//			source: func() string {
		//				return `CAPTURE cmd:'/bin/echo -n "HELLO:WORLD"'`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO:WORLD" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE double-quoted named param with embedded colon",
		//			source: func() string {
		//				return `CAPTURE cmd:"/bin/echo -n 'HELLO:WORLD'"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected CLI string %s", cmd.GetCmdString())
		//				}
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/echo" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %d", len(cliArgs))
		//				}
		//				if cliArgs[0] != "-n" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				if cliArgs[1] != "HELLO:WORLD" {
		//					return fmt.Errorf("CAPTURE has unexpected cli args: %#v", cliArgs)
		//				}
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE unquoted named param with multiple embedded colons",
		//			source: func() string {
		//				return "CAPTURE cmd:/bin/date:time:"
		//			},
		//			script: func(s *Script) error {
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if cliCmd != "/bin/date:time:" {
		//					return fmt.Errorf("CAPTURE parsed unexpected command name: %s", cliCmd)
		//				}
		//				if len(cliArgs) != 0 {
		//					return fmt.Errorf("CAPTURE parsed unexpected command args: %d", len(cliArgs))
		//				}
		//
		//				return nil
		//			},
		//		},
		//		{
		//			name: "CAPTURE with shell and quoted subproc with embedded colon",
		//			source: func() string {
		//				return `CAPTURE shell:"/bin/bash -c" cmd:"echo 'HELLO:WORLD'"`
		//			},
		//			script: func(s *Script) error {
		//				if len(s.Actions) != 1 {
		//					return fmt.Errorf("Script has unexpected actions, needs %d", len(s.Actions))
		//				}
		//				cmd := s.Actions[0].(*CaptureCommand)
		//				if cmd.Args()["cmd"] != cmd.GetCmdString() {
		//					return fmt.Errorf("CAPTURE action with unexpected command string %s", cmd.GetCmdString())
		//				}
		//				if cmd.Args()["shell"] != cmd.GetCmdShell() {
		//					return fmt.Errorf("CAPTURE action with unexpected shell %s", cmd.GetCmdShell())
		//				}
		//
		//				cliCmd, cliArgs, err := cmd.GetParsedCmd()
		//				if err != nil {
		//					return fmt.Errorf("CAPTURE command parse failed: %s", err)
		//				}
		//				if len(cliArgs) != 2 {
		//					return fmt.Errorf("CAPTURE unexpected command args parsed: %#v", cliArgs)
		//				}
		//				if cliCmd != "/bin/bash" {
		//					return fmt.Errorf("CAPTURE unexpected command parsed: %#v", cliCmd)
		//				}
		//				if cliArgs[0] != "-c" {
		//					return fmt.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
		//				}
		//				if cliArgs[1] != "echo 'HELLO:WORLD'" {
		//					return fmt.Errorf("CAPTURE has unexpected shell argument: expecting -c, got %s", cliArgs[0])
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

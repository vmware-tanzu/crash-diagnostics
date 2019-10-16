// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"strings"
	"testing"
)

type commandTest struct {
	name       string
	source     func() string
	script     func(*Script) error
	shouldFail bool
}

func runCommandTest(t *testing.T, test commandTest) {
	script, err := Parse(strings.NewReader(test.source()))
	if err != nil {
		if !test.shouldFail {
			t.Fatal(err)
		}
		t.Log(err)
		return
	}
	if test.script != nil {
		if err := test.script(script); err != nil {
			if !test.shouldFail {
				t.Fatal(err)
			}
			t.Log(err)
		}
	}
}
func TestCommandParse(t *testing.T) {
	tests := []commandTest{
		{
			name: "Preambles only",
			source: func() string {
				return "FROM local \n WORKDIR /a/b/c \n ENV a=b \n ENV c=d"
			},
			script: func(s *Script) error {
				fromCmds := s.Preambles[CmdFrom]
				if len(fromCmds) != 1 {
					return fmt.Errorf("Script has unexpected preamble %s", CmdFrom)
				}
				wdCmds := s.Preambles[CmdWorkDir]
				if len(wdCmds) != 1 {
					return fmt.Errorf("Script has  unexpected preamble %s", CmdWorkDir)
				}
				envCmds := s.Preambles[CmdEnv]
				if len(envCmds) != 2 {
					return fmt.Errorf("Script has unexpected preamble %s", envCmds)
				}
				asCmds := s.Preambles[CmdAs]
				if len(asCmds) != 1 {
					return fmt.Errorf("Script missing default preamble %s", CmdAs)
				}
				return nil
			},
		},
		{
			name: "Actions only",
			source: func() string {
				return "CAPTURE /a/b c d\n CAPTURE e f\n COPY f/g h/i/k"
			},
			script: func(s *Script) error {
				fromCmds := s.Preambles[CmdFrom]
				if len(fromCmds) != 1 {
					return fmt.Errorf("Script missing default preamble %s", CmdFrom)
				}
				wdCmds := s.Preambles[CmdWorkDir]
				if len(wdCmds) != 1 {
					return fmt.Errorf("Script missing default preamble %s", CmdWorkDir)
				}
				asCmds := s.Preambles[CmdAs]
				if len(asCmds) != 1 {
					return fmt.Errorf("Script missing preamble %s", asCmds)
				}
				actions := s.Actions
				if len(actions) != 3 {
					return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
				}
				return nil
			},
		},
		{
			name: "Preambles and actions",
			source: func() string {
				return "CAPTURE /a/b c d\n CAPTURE e f\n COPY f/g h/i/k\nWORKDIR l/m/n"
			},
			script: func(s *Script) error {
				fromCmds := s.Preambles[CmdFrom]
				if len(fromCmds) != 1 {
					return fmt.Errorf("Script missing default preamble %s", CmdFrom)
				}
				wdCmds := s.Preambles[CmdWorkDir]
				if len(wdCmds) != 1 {
					return fmt.Errorf("Script missing default preamble %s", CmdWorkDir)
				}
				dir := wdCmds[0].(*WorkdirCommand).Path()
				if dir != "l/m/n" {
					return fmt.Errorf("Script instruction WORKDIR has unexpected Dir %s", dir)
				}
				asCmds := s.Preambles[CmdAs]
				if len(asCmds) != 1 {
					return fmt.Errorf("Script missing preamble %s", asCmds)
				}
				actions := s.Actions
				if len(actions) != 3 {
					return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
				}
				return nil
			},
		},
		{
			name: "Script with comments",
			source: func() string {
				return "CAPTURE /a/b c d\n#this is a comment\n COPY f/g h/i/k\nWORKDIR l/m/n"
			},
			script: func(s *Script) error {
				actions := s.Actions
				if len(actions) != 2 {
					return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
				}
				cpCmd := s.Actions[1].(*CopyCommand)
				if len(cpCmd.Paths()) != 2 {
					return fmt.Errorf("Unexpected arg count %d for COPY in script with comment", len(cpCmd.Paths()))
				}
				return nil
			},
		},
		{
			name: "Script with only comments",
			source: func() string {
				return "#Comment line 1\n#this is a comment line 2\n # Comment line 3"
			},
			script: func(s *Script) error {
				actions := s.Actions
				if len(actions) != 0 {
					return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
				}
				preambles := s.Preambles
				if len(preambles) != 5 {
					return fmt.Errorf("Script has unexpected number of preambles %d", len(preambles))
				}
				return nil
			},
		},
		{
			name: "Script with bad preamble",
			source: func() string {
				return "CAPTURE /a/b c d\n CAPTURE e f\n COPY f/g h/i/k\nENV a|b"
			},
			shouldFail: true,
		},
		{
			name: "Script with bad action",
			source: func() string {
				return "CAPTURE\n CAPTURE e f\n COPY f/g h/i/k\nENV a|b"
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

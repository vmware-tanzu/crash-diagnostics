// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestMain(m *testing.M) {
	testcrashd.Init()
	os.Exit(m.Run())
}

type parserTest struct {
	name       string
	source     func(*testing.T) string
	script     func(*testing.T, *script.Script)
	shouldFail bool
}

func runParserTest(t *testing.T, test parserTest) {
	scr, err := Parse(strings.NewReader(test.source(t)))
	if err != nil {
		t.Fatal(err)
	}
	if test.script != nil {
		test.script(t, scr)
	}
}

func TestCommandParse(t *testing.T) {
	tests := []parserTest{
		//{
		//	name: "Preambles only",
		//	source: func(t *testing.T) string {
		//		return "FROM local \n WORKDIR /a/b/c \n ENV a=b \n ENV c=d"
		//	},
		//	script: func(t *testing.T, s *script.Script) {
		//		fromCmds := s.Preambles[script.CmdFrom]
		//		if len(fromCmds) != 1 {
		//			t.Errorf("Script has unexpected preamble %s", script.CmdFrom)
		//		}
		//		wdCmds := s.Preambles[script.CmdWorkDir]
		//		if len(wdCmds) != 1 {
		//			t.Errorf("Script has  unexpected preamble %s", script.CmdWorkDir)
		//		}
		//		envCmds := s.Preambles[script.CmdEnv]
		//		if len(envCmds) != 2 {
		//			t.Errorf("Script has unexpected preamble %s", envCmds)
		//		}
		//		asCmds := s.Preambles[script.CmdAs]
		//		if len(asCmds) != 1 {
		//			t.Errorf("Script missing default preamble %s", script.CmdAs)
		//		}
		//	},
		//},
		//{
		//	name: "Actions only",
		//	source: func() string {
		//		return "CAPTURE /a/b c d\n CAPTURE e f\n COPY f/g h/i/k"
		//	},
		//	script: func(s *Script) error {
		//		fromCmds := s.Preambles[CmdFrom]
		//		if len(fromCmds) != 1 {
		//			return fmt.Errorf("Script missing default preamble %s", CmdFrom)
		//		}
		//		wdCmds := s.Preambles[CmdWorkDir]
		//		if len(wdCmds) != 1 {
		//			return fmt.Errorf("Script missing default preamble %s", CmdWorkDir)
		//		}
		//		asCmds := s.Preambles[CmdAs]
		//		if len(asCmds) != 1 {
		//			return fmt.Errorf("Script missing preamble %s", asCmds)
		//		}
		//		actions := s.Actions
		//		if len(actions) != 3 {
		//			return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
		//		}
		//		return nil
		//	},
		//},
		//{
		//	name: "Preambles and actions",
		//	source: func() string {
		//		return "CAPTURE /a/b c d\n CAPTURE e f\n COPY f/g h/i/k\nWORKDIR l/m/n"
		//	},
		//	script: func(s *Script) error {
		//		fromCmds := s.Preambles[CmdFrom]
		//		if len(fromCmds) != 1 {
		//			return fmt.Errorf("Script missing default preamble %s", CmdFrom)
		//		}
		//		wdCmds := s.Preambles[CmdWorkDir]
		//		if len(wdCmds) != 1 {
		//			return fmt.Errorf("Script missing default preamble %s", CmdWorkDir)
		//		}
		//		dir := wdCmds[0].(*WorkdirCommand).Path()
		//		if dir != "l/m/n" {
		//			return fmt.Errorf("Script instruction WORKDIR has unexpected Dir %s", dir)
		//		}
		//		asCmds := s.Preambles[CmdAs]
		//		if len(asCmds) != 1 {
		//			return fmt.Errorf("Script missing preamble %s", asCmds)
		//		}
		//		actions := s.Actions
		//		if len(actions) != 3 {
		//			return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
		//		}
		//		return nil
		//	},
		//},
		//{
		//	name: "Script with comments",
		//	source: func() string {
		//		return "CAPTURE /a/b c d\n#this is a comment\n COPY f/g h/i/k\nWORKDIR l/m/n"
		//	},
		//	script: func(s *Script) error {
		//		actions := s.Actions
		//		if len(actions) != 2 {
		//			return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
		//		}
		//		cpCmd := s.Actions[1].(*CopyCommand)
		//		if len(cpCmd.Paths()) != 2 {
		//			return fmt.Errorf("Unexpected arg count %d for COPY in script with comment", len(cpCmd.Paths()))
		//		}
		//		return nil
		//	},
		//},
		//{
		//	name: "Script with only comments",
		//	source: func() string {
		//		return "#Comment line 1\n#this is a comment line 2\n # Comment line 3"
		//	},
		//	script: func(s *Script) error {
		//		actions := s.Actions
		//		if len(actions) != 0 {
		//			return fmt.Errorf("Script has unexpected number of actions %d", len(actions))
		//		}
		//		preambles := s.Preambles
		//		if len(preambles) != 6 {
		//			return fmt.Errorf("Script has unexpected number of preambles %d", len(preambles))
		//		}
		//		return nil
		//	},
		//},
		//{
		//	name: "Script with bad preamble",
		//	source: func() string {
		//		return "CAPTURE /a/b c d\n CAPTURE e f\n COPY f/g h/i/k\nENV a|b"
		//	},
		//	shouldFail: true,
		//},
		//{
		//	name: "Script with bad action",
		//	source: func() string {
		//		return "CAPTURE\n CAPTURE e f\n COPY f/g h/i/k\nENV a|b"
		//	},
		//	shouldFail: true,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runParserTest(t, test)
		})
	}
}

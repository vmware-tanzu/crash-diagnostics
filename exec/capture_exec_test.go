// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestExecCAPTURE(t *testing.T) {
	tests := []execTest{
		{
			name: "CAPTURE single remote command",
			source: func() string {
				src := fmt.Sprintf(`FROM 127.0.0.1:%s
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				CAPTURE /bin/echo "HELLO WORLD"`, testcrashd.DefaultSSHPort())
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				capCmd := s.Actions[0].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), fmt.Sprintf("%s.txt", sanitizeStr(capCmd.GetCmdString())))
				if _, err := os.Stat(fileName); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "CAPTURE multiple commands",
			source: func() string {
				src := `FROM 127.0.0.1:2222
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				CAPTURE /bin/echo HELLO!
				CAPTURE ls /tmp`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				cmd0 := s.Actions[0].(*script.CaptureCommand)
				cmd1 := s.Actions[1].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fname0 := filepath.Join(workdir.Path(), sanitizeStr(machine), fmt.Sprintf("%s.txt", sanitizeStr(cmd0.GetCmdString())))
				fname1 := filepath.Join(workdir.Path(), sanitizeStr(machine), fmt.Sprintf("%s.txt", sanitizeStr(cmd1.GetCmdString())))
				if _, err := os.Stat(fname0); err != nil {
					return err
				}
				if _, err := os.Stat(fname1); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "CAPTURE remote command AS user",
			source: func() string {
				src := `FROM 127.0.0.1:2222
				AS userid:${USER}
				AUTHCONFIG private-key:${HOME}/.ssh/id_rsa
				CAPTURE /bin/echo "HELLO WORLD"`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				capCmd := s.Actions[0].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), fmt.Sprintf("%s.txt", sanitizeStr(capCmd.GetCmdString())))
				if _, err := os.Stat(fileName); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "CAPTURE unquoted default with quoted subcommand",
			source: func() string {
				return `
				FROM 127.0.0.1:2222
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				CAPTURE /bin/bash -c 'echo "Hello to the World!"'`
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				capCmd := s.Actions[0].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), fmt.Sprintf("%s.txt", sanitizeStr(capCmd.GetCmdString())))
				if _, err := os.Stat(fileName); err != nil {
					return err
				}
				content, err := ioutil.ReadFile(fileName)
				if err != nil {
					return err
				}
				if strings.TrimSpace(string(content)) != "Hello to the World!" {
					return fmt.Errorf("CAPTURE generated unexpected file content: %s", content)
				}
				return nil
			},
		},
		{
			name: "CAPTURE remote command AS bad user",
			source: func() string {
				src := `FROM 127.0.0.1:2222
				AS userid:foo
				AUTHCONFIG private-key:${HOME}/.ssh/id_rsa
				CAPTURE /bin/echo "HELLO WORLD"`
				return src
			},
			exec: func(s *script.Script) error {
				e := New(s)
				return e.Execute()
			},
			shouldFail: true,
		},
		{
			name: "CAPTURE with echo on",
			source: func() string {
				src := `FROM 127.0.0.1:2222
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				CAPTURE cmd:'/bin/echo "HELLO WORLD"' echo:"on"`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				capCmd := s.Actions[0].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), fmt.Sprintf("%s.txt", sanitizeStr(capCmd.GetCmdString())))
				if _, err := os.Stat(fileName); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "CAPTURE remote command with bad AUTHCONFIG user",
			source: func() string {
				src := `FROM 127.0.0.1:22
				AUTHCONFIG username:_foouser private-key:$HOME/.ssh/id_rsa
				CAPTURE /bin/echo "HELLO WORLD"`
				return src
			},
			exec: func(s *script.Script) error {
				e := New(s)
				return e.Execute()
			},
			shouldFail: true,
		},
		{
			name: "CAPTURE bad remote command",
			source: func() string {
				src := `FROM 127.0.0.1:2222
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				CAPTURE _foo_ _bar_`
				return src
			},
			exec: func(s *script.Script) error {
				e := New(s)
				return e.Execute()
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

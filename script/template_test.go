// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestApplyTemplate(t *testing.T) {
	tests := []struct {
		name       string
		source     func() string
		compare    func(string) error
		shouldFail bool
	}{
		{
			name:   "Retrieve Home var",
			source: func() string { return "{{.Home}}/a/b/c" },
			compare: func(result string) error {
				dir, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				expected := filepath.Join(dir, "/a/b/c")
				if expected != result {
					return fmt.Errorf("Unexpected templated result, expecting [%s] got [%s]: ", expected, result)
				}
				return nil
			},
		},

		{
			name:   "Retrieve Username var",
			source: func() string { return "{{.Username}}/a/b/c" },
			compare: func(result string) error {
				usr, err := user.Current()
				if err != nil {
					return err
				}
				expected := filepath.Join(usr.Username, "/a/b/c")
				if expected != result {
					return fmt.Errorf("Unexpected templated result, expecting [%s] got [%s]: ", expected, result)
				}
				return nil
			},
		},
		{
			name:   "Retrieve Pwd var",
			source: func() string { return "{{.Pwd}}/a/b/c" },
			compare: func(result string) error {
				dir, err := os.Getwd()
				if err != nil {
					return err
				}
				expected := filepath.Join(dir, "/a/b/c")
				if expected != result {
					return fmt.Errorf("Unexpected templated result, expecting [%s] got [%s]: ", expected, result)
				}
				return nil
			},
		},
		{
			name:       "Bad or missing vars",
			source:     func() string { return "{{.Pwdi}}/a/b/c" },
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := new(bytes.Buffer)
			if err := applyTemplate(result, test.source()); err != nil {
				if !test.shouldFail {
					t.Fatal(err)
				}
				t.Log(err)
				return
			}

			if err := test.compare(result.String()); err != nil {
				if !test.shouldFail {
					t.Fatal(err)
				}
				t.Log(err)
				return
			}
		})
	}
}

func TestTemplatedScripts(t *testing.T) {
	tests := []commandTest{
		{
			name: "templated script with Home",
			source: func() string {
				return "WORKDIR {{.Home}}/.script"
			},
			script: func(s *Script) error {
				hdir, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				dirs := s.Preambles[CmdWorkDir]
				wdCmd := dirs[0].(*WorkdirCommand)
				expected := filepath.Join(hdir, "/.script")
				if wdCmd.Dir() != expected {
					return fmt.Errorf("Templated script failed, expecting WORKDIR %s, got %s", expected, wdCmd.Dir())
				}
				return nil
			},
		},
		{
			name: "templated script with Username",
			source: func() string {
				return "ENV USR={{.Username}}"
			},
			script: func(s *Script) error {
				usr, err := user.Current()
				if err != nil {
					return err
				}
				envs := s.Preambles[CmdEnv]
				envCmd := envs[0].(*EnvCommand)
				expected := fmt.Sprintf("USR=%s", usr.Username)
				if envCmd.Envs()[0] != expected {
					return fmt.Errorf("Templated script failed, expecting %s, got %s", expected, envCmd.Envs()[0])
				}
				return nil
			},
		},
		{
			name: "templated script with Pwd",
			source: func() string {
				return "WORKDIR {{.Pwd}}/.script"
			},
			script: func(s *Script) error {
				pwd, err := os.Getwd()
				if err != nil {
					return err
				}
				dirs := s.Preambles[CmdWorkDir]
				wdCmd := dirs[0].(*WorkdirCommand)
				expected := filepath.Join(pwd, "/.script")
				if wdCmd.Dir() != expected {
					return fmt.Errorf("Templated script failed, expecting WORKDIR %s, got %s", expected, wdCmd.Dir())
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

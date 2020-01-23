// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecRUN(t *testing.T) {
	tests := []execTest{
		{
			name: "RUN single command",
			source: func() string {
				return fmt.Sprintf(`FROM 127.0.0.1:%s
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN /bin/echo "HELLO WORLD"
				`, testSSHPort)
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				exitcode := os.Getenv("CMD_EXITCODE")
				if exitcode != "0" {
					return fmt.Errorf("RUN has unexpected exit code %s", exitcode)
				}

				result := os.Getenv("CMD_RESULT")
				if result != "HELLO WORLD" {
					return fmt.Errorf("RUN has unexpected CMD_RESULT: %s", result)
				}
				return nil
			},
		},
		{
			name: "RUN multiple commands",
			source: func() string {
				return fmt.Sprintf(`
				FROM 127.0.0.1:%s
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN "/bin/echo 'HELLO WORLD'"
				RUN "/bin/echo 'FROM SPACE'"
				`, testSSHPort)
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				exitcode := os.Getenv("CMD_EXITCODE")
				if exitcode != "0" {
					return fmt.Errorf("RUN has unexpected exit code %s", exitcode)
				}

				result := os.Getenv("CMD_RESULT")
				if result != "FROM SPACE" {
					return fmt.Errorf("RUN has unexpected CMD_RESULT: %s", result)
				}
				return nil
			},
		},
		{
			name: "RUN chain command result",
			source: func() string {
				return fmt.Sprintf(`
				FROM 127.0.0.1:%s
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN "/bin/echo 'HELLO WORLD'"
				RUN "/bin/echo '${CMD_RESULT} ALL'"
				`, testSSHPort)
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				exitcode := os.Getenv("CMD_EXITCODE")
				if exitcode != "0" {
					return fmt.Errorf("RUN has unexpected exit code %s", exitcode)
				}

				result := os.Getenv("CMD_RESULT")
				if result != "HELLO WORLD ALL" {
					return fmt.Errorf("RUN has unexpected CMD_RESULT: %s", result)
				}
				return nil
			},
		},
		{
			name: "RUN default param with quoted subcommand",
			source: func() string {
				return fmt.Sprintf(`
				FROM 127.0.0.1:%s
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN /bin/bash -c 'echo "Hello World"'`, testSSHPort)
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				exitcode := os.Getenv("CMD_EXITCODE")
				if exitcode != "0" {
					return fmt.Errorf("RUN has unexpected exit code %s", exitcode)
				}

				result := os.Getenv("CMD_RESULT")
				if strings.TrimSpace(result) != "Hello World" {
					return fmt.Errorf("RUN has unexpected CMD_RESULT: %s", result)
				}
				return nil
			},
		},
		{
			name: "RUN with shell and wrapped quoted subcommand",
			source: func() string {
				return fmt.Sprintf(`
				FROM 127.0.0.1:%s
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN shell:"/bin/bash -c" cmd:'echo "Hello World"'`, testSSHPort)
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				exitcode := os.Getenv("CMD_EXITCODE")
				if exitcode != "0" {
					return fmt.Errorf("RUN has unexpected exit code %s", exitcode)
				}

				result := os.Getenv("CMD_RESULT")
				if strings.TrimSpace(result) != "Hello World" {
					return fmt.Errorf("RUN has unexpected CMD_RESULT: %s", result)
				}
				return nil
			},
		},
		{
			name: "RUN with echo on",
			source: func() string {
				return fmt.Sprintf(`FROM 127.0.0.1:%s
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN cmd:'/bin/echo "HELLO WORLD"' echo:"on"
				`, testSSHPort)
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				result := os.Getenv("CMD_RESULT")
				if result != "HELLO WORLD" {
					return fmt.Errorf("RUN has unexpected CMD_RESULT: %s", result)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

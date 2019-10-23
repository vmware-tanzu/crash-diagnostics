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

func TestExecLocalRUN(t *testing.T) {
	tests := []execTest{
		{
			name: "RUN single command",
			source: func() string {
				return `RUN "/bin/echo 'HELLO WORLD'"`
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				pid := os.Getenv("CMD_PID")
				if pid == "" {
					return fmt.Errorf("RUN has unexpected pid %s", pid)
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
				return `
				RUN "/bin/echo 'HELLO WORLD'"
				RUN "/bin/echo 'FROM SPACE'"
				`
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				pid := os.Getenv("CMD_PID")
				if pid == "" {
					return fmt.Errorf("RUN has unexpected pid %s", pid)
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
				return `
				RUN "/bin/echo 'HELLO WORLD'"
				RUN "/bin/echo '${CMD_RESULT} ALL'"
				`
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				pid := os.Getenv("CMD_PID")
				if pid == "" {
					return fmt.Errorf("RUN has unexpected pid %s", pid)
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
			name: "RUN with error code",
			source: func() string {
				return `
				RUN "/bin/date --foo"
				`
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					t.Log(err)
				}
				exitcode := os.Getenv("CMD_EXITCODE")
				if exitcode == "0" {
					return fmt.Errorf("RUN has unexpected exit code %s", exitcode)
				}

				result := os.Getenv("CMD_SUCCESS")
				if result == "true" {
					return fmt.Errorf("RUN has unexpected CMD_SUCCESS: %s", result)
				}
				return nil
			},
		},
		{
			name: "RUN unquoted default with quoted subcommand",
			source: func() string {
				return `RUN /bin/bash -c 'echo "Hello World"'`
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				pid := os.Getenv("CMD_PID")
				if pid == "" {
					return fmt.Errorf("RUN has unexpected pid %s", pid)
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
				return `RUN shell:"/bin/bash -c" cmd:'echo "Hello World"'`
			},
			exec: func(s *script.Script) error {

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				pid := os.Getenv("CMD_PID")
				if pid == "" {
					return fmt.Errorf("RUN has unexpected pid %s", pid)
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

// TestExecRemoteRUN test executes a remote command using an SSH backend.
// It assumes running account has $HOME/.ssh/id_rsa private key and
// that the remote machine has public key in authorized_keys.
// If setup properly, comment out t.Skip()
func TestExecRemoteRUN(t *testing.T) {
	t.Skip(`Skipping: test requires an ssh daemon running and a
		passwordless setup using private key specified with AUTHCONFIG command`)

	tests := []execTest{
		{
			name: "RUN single command",
			source: func() string {
				return `FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN /bin/echo "HELLO WORLD"
				`
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
				return `
				FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN "/bin/echo 'HELLO WORLD'"
				RUN "/bin/echo 'FROM SPACE'"
				`
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
				return `
				FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN "/bin/echo 'HELLO WORLD'"
				RUN "/bin/echo '${CMD_RESULT} ALL'"
				`
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
				return `
				FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN /bin/bash -c 'echo "Hello World"'`
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
				return `
				FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				RUN shell:"/bin/bash -c" cmd:'echo "Hello World"'`
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

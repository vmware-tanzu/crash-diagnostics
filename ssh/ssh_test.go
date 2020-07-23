// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"strings"
	"testing"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestRun(t *testing.T) {
	usr := testcrashd.GetSSHUsername()
	pkPath := testcrashd.GetSSHPrivateKey()

	tests := []struct {
		name   string
		args   SSHArgs
		cmd    string
		result string
	}{
		{
			name:   "simple cmd",
			args:   SSHArgs{User: usr, PrivateKeyPath: pkPath, Host: "127.0.0.1", Port: testSSHPort, MaxRetries: testMaxRetries},
			cmd:    "echo 'Hello World!'",
			result: "Hello World!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected, err := Run(test.args, test.cmd)
			if err != nil {
				t.Fatal(err)
			}
			if test.result != expected {
				t.Fatalf("unexpected result %s", expected)
			}
		})
	}
}

func TestRunRead(t *testing.T) {
	usr := testcrashd.GetSSHUsername()
	pkPath := testcrashd.GetSSHPrivateKey()

	tests := []struct {
		name   string
		args   SSHArgs
		cmd    string
		result string
	}{
		{
			name:   "simple cmd",
			args:   SSHArgs{User: usr, PrivateKeyPath: pkPath, Host: "127.0.0.1", Port: testSSHPort, MaxRetries: testMaxRetries},
			cmd:    "echo 'Hello World!'",
			result: "Hello World!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader, err := RunRead(test.args, test.cmd)
			if err != nil {
				t.Fatal(err)
			}
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(reader); err != nil {
				t.Fatal(err)
			}
			expected := strings.TrimSpace(buf.String())
			if test.result != expected {
				t.Fatalf("unexpected result %s", expected)
			}
		})
	}
}

func TestSSHRunMakeCmdStr(t *testing.T) {
	tests := []struct {
		name       string
		args       SSHArgs
		cmdStr     string
		shouldFail bool
	}{
		{
			name:   "user and host",
			args:   SSHArgs{User: "sshuser", Host: "local.host"},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -p 22 sshuser@local.host",
		},
		{
			name:   "user host and pkpath",
			args:   SSHArgs{User: "sshuser", Host: "local.host", PrivateKeyPath: "/pk/path"},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -i /pk/path -p 22 sshuser@local.host",
		},
		{
			name:   "user host pkpath and proxy",
			args:   SSHArgs{User: "sshuser", Host: "local.host", PrivateKeyPath: "/pk/path", ProxyJump: &ProxyJumpArgs{User: "juser", Host: "jhost"}},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -i /pk/path -p 22 -J juser@jhost sshuser@local.host",
		},
		{
			name:       "missing host",
			args:       SSHArgs{User: "sshuser"},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := makeSSHCmdStr("ssh", test.args)
			if err != nil && !test.shouldFail {
				t.Fatal(err)
			}
			cmdFields := strings.Fields(test.cmdStr)
			resultFields := strings.Fields(result)

			for i := range cmdFields {
				if cmdFields[i] != resultFields[i] {
					t.Fatalf("unexpected command string element: %s vs. %s", cmdFields, resultFields)
				}
			}
		})
	}
}

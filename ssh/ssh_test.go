// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name   string
		args   SSHArgs
		cmd    string
		result string
	}{
		{
			name:   "simple cmd",
			args:   testSSHArgs,
			cmd:    "echo 'Hello World!'",
			result: "Hello World!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected, err := Run(test.args, nil, test.cmd)
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
	tests := []struct {
		name   string
		args   SSHArgs
		cmd    string
		result string
	}{
		{
			name:   "simple cmd",
			args:   testSSHArgs,
			cmd:    "echo 'Hello World!'",
			result: "Hello World!",
		},
		{
			name:   "simple cmd on IPv6 host",
			args:   testSSHArgsIPv6,
			cmd:    "echo 'Hello World!'",
			result: "Hello World!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader, err := RunRead(test.args, nil, test.cmd)
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
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -i /pk/path -p 22 sshuser@local.host -o \"ProxyCommand ssh -o StrictHostKeyChecking=no -W %h:%p -i /pk/path juser@jhost\"",
		},
		{
			name:   "user host and proxy",
			args:   SSHArgs{User: "sshuser", Host: "local.host", ProxyJump: &ProxyJumpArgs{User: "juser", Host: "jhost"}},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -p 22 sshuser@local.host -o \"ProxyCommand ssh -o StrictHostKeyChecking=no -W %h:%p juser@jhost\"",
		},
		{
			name:       "missing host",
			args:       SSHArgs{User: "sshuser"},
			shouldFail: true,
		},
		{
			name:   "user and IPv6 host",
			args:   SSHArgs{User: "sshuser", Host: "b::1"},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -p 22 sshuser@b::1",
		},
		{
			name:   "user IPv6 host and pkpath",
			args:   SSHArgs{User: "sshuser", Host: "b::1", PrivateKeyPath: "/pk/path"},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -i /pk/path -p 22 sshuser@b::1",
		},
		{
			name:   "user IPv6 host pkpath and proxy",
			args:   SSHArgs{User: "sshuser", Host: "b::1", PrivateKeyPath: "/pk/path", ProxyJump: &ProxyJumpArgs{User: "juser", Host: "jhost"}},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -i /pk/path -p 22 sshuser@b::1 -o \"ProxyCommand ssh -o StrictHostKeyChecking=no -W %h:%p -i /pk/path juser@jhost\"",
		},
		{
			name:   "user IPv6 host and IPv6 proxy",
			args:   SSHArgs{User: "sshuser", Host: "b::1", ProxyJump: &ProxyJumpArgs{User: "juser", Host: "::a"}},
			cmdStr: "ssh -q -o StrictHostKeyChecking=no -p 22 sshuser@b::1 -o \"ProxyCommand ssh -o StrictHostKeyChecking=no -W %h:%p juser@::a\"",
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

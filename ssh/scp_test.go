// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	pkPath := sshSvr.PrivateKey()
	sshArgs := SSHArgs{User: testSSHUsername, PrivateKeyPath: pkPath, Host: "127.0.0.1", Port: testSSHPort, MaxRetries: testMaxRetries}

	tests := []struct {
		name        string
		sshArgs     SSHArgs
		rootDir     string
		remoteFiles map[string]string
		srcFile     string
		fileContent string
	}{
		{
			name:        "copy single file",
			sshArgs:     sshArgs,
			rootDir:     "/tmp/crashd",
			remoteFiles: map[string]string{"foo.txt": "FooBar"},
			srcFile:     "foo.txt",
			fileContent: "FooBar",
		},
		{
			name:        "copy single file in dir",
			sshArgs:     sshArgs,
			rootDir:     "/tmp/crashd",
			remoteFiles: map[string]string{"foo/bar.txt": "FooBar"},
			srcFile:     "foo/bar.txt",
			fileContent: "FooBar",
		},
		{
			name:        "copy dir",
			sshArgs:     sshArgs,
			rootDir:     "/tmp/crashd",
			remoteFiles: map[string]string{"bar/foo.csv": "FooBar", "bar/bar.txt": "BarBar"},
			srcFile:     "bar/",
			fileContent: "FooBar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				for file, _ := range test.remoteFiles {
					RemoveTestSSHFile(t, test.sshArgs, file)
				}

				if err := os.RemoveAll(test.rootDir); err != nil {
					t.Fatal(err)
				}
			}()

			// setup remote files
			for file, content := range test.remoteFiles {
				MakeTestSSHFile(t, test.sshArgs, file, content)
			}

			if err := CopyFrom(test.sshArgs, test.rootDir, test.srcFile); err != nil {
				t.Fatal(err)
			}

			expectedPath := filepath.Join(test.rootDir, test.srcFile)
			finfo, err := os.Stat(expectedPath)
			if err != nil {
				t.Fatal(err)
			}

			if finfo.IsDir() {
				finfos, err := ioutil.ReadDir(expectedPath)
				if err != nil {
					t.Fatal(err)
				}
				if len(finfos) < len(test.remoteFiles) {
					t.Errorf("expecting %d copied files, got %d", len(finfos), len(test.remoteFiles))
				}
			} else {
				if getTestFileContent(t, expectedPath) != test.fileContent {
					t.Error("unexpected file content")
				}
			}

		})
	}
}

func TestMakeSCPCmdStr(t *testing.T) {
	tests := []struct {
		name       string
		args       SSHArgs
		cmdStr     string
		source     string
		shouldFail bool
	}{
		{
			name:   "user and host",
			args:   SSHArgs{User: "sshuser", Host: "local.host"},
			source: "/tmp/any",
			cmdStr: "scp -rpq -o StrictHostKeyChecking=no -P 22 sshuser@local.host:/tmp/any",
		},
		{
			name:   "user host and pkpath",
			args:   SSHArgs{User: "sshuser", Host: "local.host", PrivateKeyPath: "/pk/path"},
			source: "/foo/bar",
			cmdStr: "scp -rpq -o StrictHostKeyChecking=no -i /pk/path -P 22 sshuser@local.host:/foo/bar",
		},
		{
			name:   "user host pkpath and proxy",
			args:   SSHArgs{User: "sshuser", Host: "local.host", PrivateKeyPath: "/pk/path", ProxyJump: &ProxyJumpArgs{User: "juser", Host: "jhost"}},
			source: "userFile",
			cmdStr: "scp -rpq -o StrictHostKeyChecking=no -i /pk/path -P 22 -J juser@jhost sshuser@local.host:userFile",
		},
		{
			name:       "missing host",
			args:       SSHArgs{User: "sshuser"},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := makeSCPCmdStr("scp", test.args, test.source)
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

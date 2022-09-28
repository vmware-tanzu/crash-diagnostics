// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFrom(t *testing.T) {
	tests := []struct {
		name        string
		sshArgs     SSHArgs
		remoteFiles map[string]string
		srcFile     string
		fileContent string
	}{
		{
			name:        "copy single file",
			sshArgs:     testSSHArgs,
			remoteFiles: map[string]string{"foo.txt": "FooBar"},
			srcFile:     "foo.txt",
			fileContent: "FooBar",
		},
		{
			name:        "copy single file in dir",
			sshArgs:     testSSHArgs,
			remoteFiles: map[string]string{"foo/bar.txt": "FooBar"},
			srcFile:     "foo/bar.txt",
			fileContent: "FooBar",
		},
		{
			name:        "copy dir",
			sshArgs:     testSSHArgs,
			remoteFiles: map[string]string{"bar/foo.csv": "FooBar", "bar/bar.txt": "BarBar"},
			srcFile:     "bar/",
			fileContent: "FooBar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				for file := range test.remoteFiles {
					RemoveRemoteTestSSHFile(t, test.sshArgs, file)
				}
			}()

			// setup fake files
			for file, content := range test.remoteFiles {
				MakeRemoteTestSSHFile(t, test.sshArgs, file, content)
			}

			if err := CopyFrom(test.sshArgs, nil, support.TmpDirRoot(), test.srcFile); err != nil {
				t.Fatal(err)
			}

			// validate copied files/dir
			expectedPath := filepath.Join(support.TmpDirRoot(), test.srcFile)
			finfo, err := os.Stat(expectedPath)
			if err != nil {
				t.Fatal(err)
			}

			if finfo.IsDir() {
				finfos, err := os.ReadDir(expectedPath)
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

func TestCopyTo(t *testing.T) {
	tests := []struct {
		name        string
		sshArgs     SSHArgs
		localFiles  map[string]string
		file        string
		fileContent string
	}{
		{
			name:        "copy single file to remote",
			sshArgs:     testSSHArgs,
			localFiles:  map[string]string{"local-foo.txt": "FooBar"},
			file:        "local-foo.txt",
			fileContent: "FooBar",
		},
		{
			name:        "copy single file in dir to remote",
			sshArgs:     testSSHArgs,
			localFiles:  map[string]string{"local-foo/local-bar.txt": "FooBar"},
			file:        "local-foo/local-bar.txt",
			fileContent: "FooBar",
		},
		{
			name:        "copy dir entire dir to remote",
			sshArgs:     testSSHArgs,
			localFiles:  map[string]string{"local-bar/local-foo.csv": "FooBar", "local-bar/local-bar.txt": "BarBar"},
			file:        "local-bar/",
			fileContent: "FooBar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				for file := range test.localFiles {
					RemoveLocalTestFile(t, filepath.Join(support.TmpDirRoot(), file))
					RemoveRemoteTestSSHFile(t, test.sshArgs, file)
				}
			}()

			// setup fake local files
			for file, content := range test.localFiles {
				MakeLocalTestFile(t, filepath.Join(support.TmpDirRoot(), file), content)
			}

			// create remote dir if needed
			// setup remote dir if needed
			MakeRemoteTestSSHDir(t, test.sshArgs, test.file)

			sourceFile := filepath.Join(support.TmpDirRoot(), test.file)
			t.Logf("copyTo: copying %s -to-> %s", sourceFile, test.file)
			if err := CopyTo(test.sshArgs, nil, sourceFile, test.file); err != nil {
				t.Fatal(err)
			}

			// validate copied files/dir
			AssertRemoteTestSSHFile(t, test.sshArgs, test.file)

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
			name:   "default",
			args:   SSHArgs{User: "sshuser", Host: "local.host"},
			source: "/tmp/any",
			cmdStr: "scp -rpq -o StrictHostKeyChecking=no -P 22",
		},
		{
			name:   "pkpath",
			args:   SSHArgs{User: "sshuser", Host: "local.host", PrivateKeyPath: "/pk/path"},
			source: "/foo/bar",
			cmdStr: "scp -rpq -o StrictHostKeyChecking=no -i /pk/path -P 22",
		},
		{
			name:   "pkpath and proxy",
			args:   SSHArgs{User: "sshuser", Host: "local.host", PrivateKeyPath: "/pk/path", ProxyJump: &ProxyJumpArgs{User: "juser", Host: "jhost"}},
			source: "userFile",
			cmdStr: "scp -rpq -o StrictHostKeyChecking=no -i /pk/path -P 22 -J juser@jhost",
		},
		{
			name:       "missing host",
			args:       SSHArgs{User: "sshuser"},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := makeSCPCmdStr("scp", test.args)
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

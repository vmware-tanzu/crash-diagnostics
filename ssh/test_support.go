// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func MakeTestSSHDir(t *testing.T, args SSHArgs, dir string) {
	t.Logf("creating test dir over SSH: %s", dir)
	_, err := Run(args, fmt.Sprintf(`mkdir -p %s`, dir))
	if err != nil {
		t.Fatal(err)
	}
	// validate
	result, _ := Run(args, fmt.Sprintf(`ls %s`, dir))
	t.Logf("dir created: %s", result)
}

func MakeTestSSHFile(t *testing.T, args SSHArgs, fileName, content string) {
	srcDir := filepath.Dir(fileName)
	if len(srcDir) > 0 && srcDir != "." {
		MakeTestSSHDir(t, args, srcDir)
	}

	t.Logf("creating test file over SSH: %s", fileName)
	_, err := Run(args, fmt.Sprintf(`echo '%s' > %s`, content, fileName))
	if err != nil {
		t.Fatal(err)
	}

	result, _ := Run(args, fmt.Sprintf(`ls %s`, fileName))
	t.Logf("file created: %s", result)
}

func RemoveTestSSHFile(t *testing.T, args SSHArgs, fileName string) {
	t.Logf("removing test file over SSH: %s", fileName)
	_, err := Run(args, fmt.Sprintf(`rm -rf %s`, fileName))
	if err != nil {
		t.Fatal(err)
	}
}

func getTestFileContent(t *testing.T, fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(buf.String())
}

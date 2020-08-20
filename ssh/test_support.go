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

//////func mountTestSSHFile(t *testing.T, mountDir, fileName, content string) {
//////	srcDir := filepath.Dir(fileName)
//////	if len(srcDir) > 0 && srcDir != "." {
//////		mountTestSSHDir(t, mountDir, srcDir)
//////	}
//////
//////	filePath := filepath.Join(mountDir, fileName)
//////	t.Logf("mounting test file in SSH: %s", filePath)
//////	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
//////		t.Fatal(err)
//////	}
//////}
////
////func mountTestSSHDir(t *testing.T, mountDir, dir string) {
////	t.Logf("mounting dir in SSH: %s", dir)
////	mountPath := filepath.Join(mountDir, dir)
////	if err := os.MkdirAll(mountPath, 0754); err != nil && !os.IsExist(err) {
////		t.Fatal(err)
////	}
////}
//
//func removeTestSSHFile(t *testing.T, mountDir, fileName string) {
//	t.Logf("removing file mounted in SSH: %s", fileName)
//	filePath := filepath.Join(mountDir, fileName)
//	if err := os.RemoveAll(filePath); err != nil && !os.IsNotExist(err) {
//		t.Fatal(err)
//	}
//}

func makeTestSSHDir(t *testing.T, args SSHArgs, dir string) {
	t.Logf("creating test dir over SSH: %s", dir)
	_, err := Run(args, fmt.Sprintf(`mkdir -p %s`, dir))
	if err != nil {
		t.Fatal(err)
	}
	// validate
	result, _ := Run(args, fmt.Sprintf(`ls %s`, dir))
	t.Logf("dir created: %s", result)
}

func MakeTestSSHFile(t *testing.T, args SSHArgs, filePath, content string) {
	srcDir := filepath.Dir(filePath)
	if len(srcDir) > 0 && srcDir != "." {
		makeTestSSHDir(t, args, srcDir)
	}

	t.Logf("creating test file over SSH: %s", filePath)
	_, err := Run(args, fmt.Sprintf(`echo '%s' > %s`, content, filePath))
	if err != nil {
		t.Fatal(err)
	}

	result, _ := Run(args, fmt.Sprintf(`ls %s`, filePath))
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

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	strSanitization = regexp.MustCompile(`[\s\"\'/\.\:]`)
)

func sanitizeStr(cmd string) string {
	return strSanitization.ReplaceAllString(cmd, "_")
}

func writeFile(writer io.Writer, source io.Reader) error {
	if _, err := io.Copy(writer, source); err != nil {
		return err
	}
	return nil
}

func writeError(writer io.Writer, err error) error {
	errReader := strings.NewReader(err.Error())
	return writeFile(writer, errReader)
}

func getFileForCaptureCmd(cmdStr, workdir, output string) (*os.File, error) {
	var outfile *os.File
	switch output {
	case OutputStdout:
		outfile = os.Stdout
		logrus.Debugf("Routing result for [%s] to stdout", cmdStr)
	case OutputStderr:
		outfile = os.Stderr
		logrus.Debugf("Routing result for [%s] to stderr", cmdStr)
	default:
		fileName := fmt.Sprintf("%s.txt", sanitizeStr(cmdStr))
		filePath := filepath.Join(workdir, fileName)
		logrus.Debugf("Creating file %s for [%s]", filePath, cmdStr)
		f, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		outfile = f
	}
	return outfile, nil
}

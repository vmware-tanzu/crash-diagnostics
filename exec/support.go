// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io"
	"os"
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

func writeCmdOutput(source io.Reader, filePath string, echo bool, cmd string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := source
	if echo {
		fmt.Fprintf(file, "%s\n", cmd)
		fmt.Fprintf(os.Stdout, "%s\n", cmd)

		reader = io.TeeReader(source, os.Stdout)
	}

	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	logrus.Debugf("Wrote file %s", filePath)

	return nil
}

func writeCmdError(err error, filePath string, cmdStr string) error {
	errReader := strings.NewReader(err.Error())
	return writeCmdOutput(errReader, filePath, false, cmdStr)
}

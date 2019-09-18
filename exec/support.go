// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
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

func writeFile(source io.Reader, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, source); err != nil {
		return err
	}
	logrus.Debugf("Wrote file %s", filePath)

	return nil
}

func writeError(err error, filePath string) error {
	errReader := strings.NewReader(err.Error())
	return writeFile(errReader, filePath)
}

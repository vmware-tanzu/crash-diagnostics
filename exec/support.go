// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"io"
	"regexp"
	"strings"
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

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package functions

import (
	"os"
	"path/filepath"
)

func MakeDir(path string, mode os.FileMode) error {
	if _, err := os.Stat(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(path, mode); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// ExpandPath converts path to include the full home dir path when prefixed with `~`.
func ExpandPath(path string) (string, error) {
	if path[0] != '~' {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, path[1:]), nil
}

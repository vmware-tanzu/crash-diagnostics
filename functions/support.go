// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package functions

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func Result(funcName FunctionName, result interface{}) (starlark.Value, error) {
	starResult := new(starlarkstruct.Struct)
	if err := typekit.Go(result).Starlark(starResult); err != nil {
		return Error(funcName, fmt.Errorf("conversion error: %v", err))
	}
	return starResult, nil
}
func Error(funcName FunctionName, err error) (starlark.Value, error) {
	return starlark.None, fmt.Errorf("%s: failed: %s", funcName, err)
}

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

func DefaultGid() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Gid
}

func DefaultUid() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Uid
}

func DefaultUsername() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Username
}
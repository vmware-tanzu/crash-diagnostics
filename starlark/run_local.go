// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"

	"github.com/vladimirvivien/echo"
	"go.starlark.net/starlark"
)

// runLocalFunc is a built-in starlark function that runs a provided command on the local machine.
// It returns the result of the command as struct containing information about the executed command.
// Starlark format: run_local(<command string>)
func runLocalFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var cmdStr string
	if args != nil && args.Len() == 1 {
		cmd, ok := args.Index(0).(starlark.String)
		if !ok {
			return starlark.None, fmt.Errorf("%s: command must be a string", identifiers.runLocal)
		}
		cmdStr = string(cmd)
	}

	p := echo.New().RunProc(cmdStr)
	if p.Err() != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.runLocal, p.Err())
	}

	return starlark.String(p.Result()), nil
}
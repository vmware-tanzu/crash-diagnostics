// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vladimirvivien/echo"
	"go.starlark.net/starlark"
)

// captureLocalFunc is a built-in starlark function that runs a provided command on the local machine.
// The output of the command is stored in a file at a specified location under the workdir directory.
// Starlark format: run_local(cmd=<command> [,workdir=path][,file_name=name][,desc=description])
func captureLocalFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var cmdStr, workdir, fileName, desc string
	if err := starlark.UnpackArgs(
		identifiers.captureLocal, args, kwargs,
		"cmd", &cmdStr,
		"workdir?", &workdir,
		"file_name?", &fileName,
		"desc?", &desc,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.captureLocal, err)
	}

	if len(workdir) == 0 {
		dir, err := getWorkdirFromThread(thread)
		if err != nil {
			return starlark.None, err
		}
		workdir = dir
	}
	if len(fileName) == 0 {
		fileName = fmt.Sprintf("%s.txt", sanitizeStr(cmdStr))
	}

	filePath := filepath.Join(workdir, fileName)
	if err := os.MkdirAll(workdir, 0744); err != nil && !os.IsExist(err) {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.captureLocal, err)
	}

	p := echo.New().RunProc(cmdStr)
	if p.Err() != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.captureLocal, p.Err())
	}

	if err := captureOutput(p.Out(), filePath, desc); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.captureLocal, err)
	}

	return starlark.String(filePath), nil
}

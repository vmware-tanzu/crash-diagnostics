// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run_local

import (
	"github.com/vladimirvivien/echo"
	"go.starlark.net/starlark"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

func (c *cmd) Run(t *starlark.Thread, args Args) Result {
	proc := echo.New().RunProc(args.Cmd)
	var err string
	if proc.Err() != nil {
		err = proc.Err().Error()
	}
	return Result{
		Error:    err,
		Pid:      int64(proc.ID()),
		Result:   proc.Result(),
		ExitCode: int64(proc.ExitCode()),
	}

}

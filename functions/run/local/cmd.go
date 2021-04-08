// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"fmt"

	"github.com/vladimirvivien/echo"
	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/run"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

func (c *cmd) Run(t *starlark.Thread, p interface{}) (functions.CommandResult, error) {
	params, ok := p.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected param type: %T", p)
	}
	proc := echo.New().RunProc(params)

	var errmsg string
	if proc.Err() != nil {
		errmsg = proc.Err().Error()
	}

	val := run.LocalProc{
		Pid:      int64(proc.ID()),
		Result:   proc.Result(),
		ExitCode: int64(proc.ExitCode()),
	}

	return functions.NewResult(val).AddError(errmsg), nil
}

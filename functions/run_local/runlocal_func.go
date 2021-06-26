// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run_local

import (
	"fmt"

	"github.com/vladimirvivien/echo"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("run_local")
	Func    = runLocalFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	builtins.Register(Name, Builtin)
}

// runLocalFunc is a built-in starlark function that runs a provided command on the local machine.
// Starlark format: result = run_local(cmd="script-command")
func runLocalFunc(thread *starlark.Thread, b *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := Run(thread, args)

	// convert and return result
	return functions.Result(Name, result)
}

// Run executes the command function
func Run(t *starlark.Thread, args Args) Result {
	proc := echo.New().RunProc(args.Cmd)
	if proc.Err() != nil {
		return Result{Error: proc.Err().Error()}
	}
	return Result{
		Proc: LocalProc{
			Pid:      int64(proc.ID()),
			Result:   proc.Result(),
			ExitCode: int64(proc.ExitCode()),
		},
	}
}

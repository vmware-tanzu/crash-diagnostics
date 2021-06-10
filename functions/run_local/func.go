// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run_local

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("run_local")
	Func    = runLocalFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// runLocalFunc is a built-in starlark function that runs a provided command on the local machine.
// Starlark format: result = run_local(cmd="script-command")
func runLocalFunc(thread *starlark.Thread, b *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := newCmd().Run(thread, args)

	// convert and return result
	return functions.Result(Name, result)
}

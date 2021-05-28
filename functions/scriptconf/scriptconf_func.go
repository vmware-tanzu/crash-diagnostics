// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

var (
	Name    = functions.FunctionName("script_config")
	Func    = scriptConfigFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	builtins.Register(Name, Builtin)
}

// scriptConfigFunc implements a starlark built-in function that gathers and stores configuration
// settings for a running script.
//
// Example:
//    script_config(workdir=path, default_shell=shellpath, requires=["command0",...,"commandN"])
//
// Args:
//   - workdir string - a path that can be used as work directory during script exec
//   - gid string - the default group id to use when executing an OS command
//   - uid string - a default userid to use when executing an OS command
//   - default_shell string - path to a shell program that can be used as default (i.e. /bin/sh)
//   - requires [] string - a list of paths for commands that should be on the machine where script is executed
//   - use_ssh_agent bool - specifies if an ssh-agent should be setup for private key management
func scriptConfigFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.FuncError(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := newCmd().Run(thread, args)

	// convert and return result
	starResult := new(starlarkstruct.Struct)
	if err := typekit.Go(result).Starlark(starResult); err != nil {
		return functions.FuncError(Name, fmt.Errorf("conversion error: %v", err))
	}
	return starResult, nil

}

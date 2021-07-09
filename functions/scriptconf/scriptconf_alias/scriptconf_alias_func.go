// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package scriptconf_alias implements an alias to built-in function make_scriptconf.
package scriptconf_alias

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf/make_scriptconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	// Func is the built-in function that implements an alias to make_scriptconf.Func.
	// This alias returns the scriptconf.Result.Config value directly as a convenience.
	// However, it will stop the script if an error occurs. For better error-handling
	// use the make_script_config function directly in scripts.
	Func       = scriptConfigFunc
	Name       = functions.FunctionName("script_config")
	Builtin    = starlark.NewBuiltin(string(Name), Func)
	Identifier = string(Name)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// scriptConfigFunc is the built-in function that implements an alias to make_scriptconf.Func.
// This alias returns the scriptconf.Result.Config value directly as a convenience.
// However, it will stop the script if an error occurs. For better error-handling
// use the make_script_config function directly in scripts.
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
//
// Returns
//  - Config: a scriptconf.Config containing the configuration data
func scriptConfigFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args scriptconf.Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := make_scriptconf.Run(thread, args)

	// return fatal error/stop script
	if result.Error != "" {
		return starlark.None, errors.New(result.Error)
	}

	// save/overwrite config result in thread
	thread.SetLocal(scriptconf.Identifier, result.Config)

	// convert and return result
	return functions.Result(Name, result.Config)
}

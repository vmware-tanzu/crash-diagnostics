// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

var (
	FuncName = "script_config"
	Func     = scriptConfigFunc
	Builtin  = starlark.NewBuiltin(FuncName, Func)
)

// Register
func init() {
	builtins.Register(FuncName, Builtin)
}

// scriptConfigFunc implements a starlark built-in function that gathers and stores configuration
// settings for a running script.
//
// Example:
//    script_config(workdir=path, default_shell=shellpath, requires=["command0",...,"commandN"])
//
// Params:
//   - workdir string - a path that can be used as work directory during script exec
//   - gid string - the default group id to use when executing an OS command
//   - uid string - a default userid to use when executing an OS command
//   - default_shell string - path to a shell program that can be used as default (i.e. /bin/sh)
//   - requires [] string - a list of paths for commands that should be on the machine where script is executed
//   - use_ssh_agent bool - specifies if an ssh-agent should be setup for private key management
func scriptConfigFunc(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var workdir, gid, uid, defaultShell string
	var useSSHAgent bool
	requires := starlark.NewList([]starlark.Value{})

	if err := starlark.UnpackArgs(
		FuncName, args, kwargs,
		"workdir?", &workdir,
		"gid?", &gid,
		"uid?", &uid,
		"default_shell?", &defaultShell,
		"requires?", &requires,
		"use_ssh_agent?", &useSSHAgent,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", FuncName, err)
	}

	var progReqs []string
	if err := typekit.Starlark(requires).Go(&progReqs); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", FuncName, err)
	}

	params := Params{
		Workdir:      workdir,
		Gid:          gid,
		Uid:          uid,
		DefaultShell: defaultShell,
		UseSSHAgent:  useSSHAgent,
		Requires:     progReqs,
	}

	result, err := newCmd().Run(thread, params)
	if err != nil {
		return starlark.None, fmt.Errorf("%s: %s", FuncName, err)
	}

	// for configuration type commands, return only config value
	var confStruct starlarkstruct.Struct
	if err := typekit.Go(result.Value()).Starlark(&confStruct); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", FuncName, err)
	}

	return &confStruct, nil
}

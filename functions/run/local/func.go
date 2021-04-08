// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

var (
	FuncName = "run_local"
	Func     = runLocalFunc
	Builtin  = starlark.NewBuiltin(FuncName, Func)
)

// runLocalFunc is a built-in starlark function that runs a provided command on the local machine.
// Starlark format: result = run_local(cmd="script-command")
func runLocalFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var cmdParam string
	if err := starlark.UnpackArgs(
		FuncName, args, kwargs,
		"cmd", &cmdParam,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", FuncName, err)
	}

	result, err := newCmd().Run(thread, cmdParam)
	if err != nil {
		return starlark.None, fmt.Errorf("%s: command failed: %s", FuncName, err)
	}

	var star starlarkstruct.Struct
	if err := typekit.Go(result.Value()).Starlark(&star); err != nil {
		return starlark.None, fmt.Errorf("%s: conversion error: %s", FuncName, err)
	}

	return &star, nil
}

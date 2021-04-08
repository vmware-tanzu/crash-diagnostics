// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"

	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

var (
	FuncName          = "archive"
	Func              = archiveFunc
	Builtin           = starlark.NewBuiltin(FuncName, Func)
	DefaultBundleName = "archive.tar.gz"
)

// Register
func init() {
	builtins.Register(FuncName, Builtin)
}

// archiveFunc implements a Starlark.Builtin function that can be used to bundle to create a
// tar file bundle.
// Script example: archive(output_file=<file name> ,source_paths=[<path list>])
func archiveFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var argOutFile string
	var argPaths *starlark.List

	if err := starlark.UnpackArgs(
		FuncName, args, kwargs,
		"output_file?", &argOutFile,
		"source_paths", &argPaths,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", FuncName, err)
	}

	var paths []string
	if err := typekit.Starlark(argPaths).Go(&paths); err != nil {
		return starlark.None, fmt.Errorf("%s: type conversion error: %s", FuncName, err)
	}

	params := Params{
		SourcePaths: paths,
		OutputFile:  argOutFile,
	}

	result, err := newCmd().Run(thread, params)
	if err != nil {
		return starlark.None, fmt.Errorf("%s: command failed: %s", FuncName, err)
	}

	return functions.MakeFuncResult(result)
}

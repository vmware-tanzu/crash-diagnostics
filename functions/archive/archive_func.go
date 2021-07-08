// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name              = functions.FunctionName("archive")
	Func              = archiveFunc
	Builtin           = starlark.NewBuiltin(string(Name), Func)
	DefaultBundleName = "archive.tar.gz"
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// archiveFunc implements a Starlark.Builtin function that can be used to bundle to create a
// tar file bundle.
// Script example: archive(output_file=<file name> ,source_paths=[<path list>])
func archiveFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var params Args
	if err := typekit.KwargsToGo(kwargs, &params); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	// execute command
	result := Run(thread, params)

	// convert and return result
	return functions.Result(Name, result)
}

// Run executes the command and returns a result
func Run(t *starlark.Thread, params Args) Result {
	if params.OutputFile == "" {
		params.OutputFile = DefaultBundleName
	}

	if len(params.SourcePaths) == 0 {
		return Result{Error: "no source path provided"}
	}

	if err := archiver.Tar(params.OutputFile, params.SourcePaths...); err != nil {
		return Result{Error: fmt.Sprintf("%s failed: %s", Name, err)}
	}

	info, err := os.Stat(params.OutputFile)
	if err != nil {
		return Result{Error: fmt.Sprintf("%s: stat failed: %s", Name, err)}
	}

	return Result{Archive: Archive{Size: info.Size(), OutputFile: params.OutputFile}}
}

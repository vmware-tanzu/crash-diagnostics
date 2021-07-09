// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package make_archive

import (
	"fmt"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/archive"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name              = functions.FunctionName("make_archive")
	Func              = makeArchiveFunc
	Builtin           = starlark.NewBuiltin(string(Name), Func)
	DefaultBundleName = "archive.tar.gz"
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// makeArchiveFunc implements a Starlark.Builtin function that can be used to bundle to create a
// tar file bundle. It returns a Result containing the Archive value or an error if one occured.
// Script example: make_archive(output_file=<file name> ,source_paths=[<path list>])
func makeArchiveFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var params archive.Args
	if err := typekit.KwargsToGo(kwargs, &params); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	// execute command
	result := Run(thread, params)

	// convert and return result
	return functions.Result(Name, result)
}

// Run executes the command and returns a result
func Run(t *starlark.Thread, params archive.Args) archive.Result {
	if params.OutputFile == "" {
		params.OutputFile = DefaultBundleName
	}

	if len(params.SourcePaths) == 0 {
		return archive.Result{Error: "no source path provided"}
	}

	if err := archiver.Tar(params.OutputFile, params.SourcePaths...); err != nil {
		return archive.Result{Error: fmt.Sprintf("%s failed: %s", Name, err)}
	}

	info, err := os.Stat(params.OutputFile)
	if err != nil {
		return archive.Result{Error: fmt.Sprintf("%s: stat failed: %s", Name, err)}
	}

	return archive.Result{Archive: archive.Archive{Size: info.Size(), OutputFile: params.OutputFile}}
}

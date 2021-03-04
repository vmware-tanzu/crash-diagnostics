// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"

	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

var (
	FuncName          = "archive"
	Func              = archiveFunc
	DefaultBundleName = "archive.tar.gz"
)

// Func implements a Starlark.Builtin function that can be used to bundle to create a
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

	if len(argOutFile) == 0 {
		argOutFile = DefaultBundleName
	}

	if len(paths) == 0 {
		return starlark.None, fmt.Errorf("%s: one or more paths required", FuncName)
	}


	if err := archiver.Tar(argOutFile, paths...); err != nil {
		return starlark.None, fmt.Errorf("%s failed: %s", FuncName, err)
	}

	return starlark.String(argOutFile), nil
}

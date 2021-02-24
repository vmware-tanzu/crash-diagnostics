// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"

	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
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
	var outputFile string
	var paths *starlark.List

	if err := starlark.UnpackArgs(
		FuncName, args, kwargs,
		"output_file?", &outputFile,
		"source_paths", &paths,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", FuncName, err)
	}

	if len(outputFile) == 0 {
		outputFile = DefaultBundleName
	}

	if paths != nil && paths.Len() == 0 {
		return starlark.None, fmt.Errorf("%s: one or more paths required", FuncName)
	}

	if err := archiver.Tar(outputFile, functions.ToStringSlice(paths)...); err != nil {
		return starlark.None, fmt.Errorf("%s failed: %s", FuncName, err)
	}

	return starlark.String(outputFile), nil
}

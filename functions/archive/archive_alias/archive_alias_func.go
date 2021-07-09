// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package archive_alias implements an starlark built-in function for
// as an alias command make_archive.  This alias fails if an error
// occurs during execution.  If error-handling is needed, use package make_archive.
package archive_alias

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/archive"
	"github.com/vmware-tanzu/crash-diagnostics/functions/archive/make_archive"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	// Func is the  starlark built-in function that implements command `archive`
	// which is an alias for command `make_archive`, implemented in package `make_arvhive`.
	// This alias returns an Archive value or a fatal error fatal error if one occurs.
	// For complete error-handling in script execution, users should use make_archive.
	Func    = archiveFunc
	Name    = functions.FunctionName("archive")
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// archiveFunc implements a starlark builtin function as an alias for make_archive.
// This alias works the same way as make_archive, however it returns an Archive value
// or returns a fatal error if there is one.
// Script example: archive(output_file=<file name> ,source_paths=[<path list>])
func archiveFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var params archive.Args
	if err := typekit.KwargsToGo(kwargs, &params); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	// execute command, if error, stop execution.
	result := make_archive.Run(thread, params)
	if result.Error != "" {
		return starlark.None, fmt.Errorf(result.Error)
	}

	// return result
	return functions.Result(Name, result.Archive)
}

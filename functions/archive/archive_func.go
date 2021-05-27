// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
)

var (
	Name              = functions.FunctionName("archive")
	Func              = archiveFunc
	Builtin           = starlark.NewBuiltin(string(Name), Func)
	DefaultBundleName = "archive.tar.gz"
)

// Register
func init() {
	builtins.Register(Name, Builtin)
}

// archiveFunc implements a Starlark.Builtin function that can be used to bundle to create a
// tar file bundle.
// Script example: archive(output_file=<file name> ,source_paths=[<path list>])
func archiveFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var params Args
	if err := typekit.KwargsToGo(kwargs, &params); err != nil {
		return functions.FuncError(Name, fmt.Errorf("%s: %s", Name, err))
	}

	// execute command
	result := newCmd().Run(thread, params)

	// convert and return result
	starResult := new(starlarkstruct.Struct)
	if err := typekit.Go(result).Starlark(starResult); err != nil {
		return functions.FuncError(Name, fmt.Errorf("conversion error: %v", err))
	}
	return starResult, nil
}

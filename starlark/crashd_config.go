// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// addDefaultCrashdConf initalizes a Starlark Dict with default
// crashd_config configuration data
func addDefaultCrashdConf(thread *starlark.Thread) error {
	args := []starlark.Tuple{
		starlark.Tuple{starlark.String("gid"), starlark.String(getGid())},
		starlark.Tuple{starlark.String("uid"), starlark.String(getUid())},
		starlark.Tuple{starlark.String("workdir"), starlark.String(defaults.workdir)},
		starlark.Tuple{starlark.String("output_path"), starlark.String(defaults.outPath)},
	}

	_, err := crashdConfigFn(thread, nil, nil, args)
	if err != nil {
		return err
	}

	return nil
}

// crashConfig is built-in starlark function that saves and returns the kwargs as a struct value.
// Starlark format: crashd_config(conf0=val0, ..., confN=ValN)
func crashdConfigFn(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var dictionary starlark.StringDict
	if kwargs != nil {
		dict, err := kwargsToStringDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		dictionary = dict
	}

	structVal := starlarkstruct.FromStringDict(starlarkstruct.Default, dictionary)

	// save values to be used as default
	thread.SetLocal(identifiers.crashdCfg, structVal)

	// return values as a struct (i.e. config.arg0, ... , config.argN)
	return starlark.None, nil
}

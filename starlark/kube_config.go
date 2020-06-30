// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// kubeConfigFn is built-in starlark function that wraps the kwargs into a dictionary value.
// The result is also added to the thread for other built-in to access.
func kubeConfigFn(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var dictionary starlark.StringDict

	if kwargs != nil {
		dict, err := kwargsToStringDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		dictionary = dict
	}
	structVal := starlarkstruct.FromStringDict(starlarkstruct.Default, dictionary)

	// save dict to be used as default
	thread.SetLocal(identifiers.kubeCfg, structVal)

	return structVal, nil
}

// addDefaultKubeConf initializes a Starlark Dict with default
// KUBECONFIG configuration data
func addDefaultKubeConf(thread *starlark.Thread) error {
	args := []starlark.Tuple{
		{starlark.String("path"), starlark.String(defaults.kubeconfig)},
	}

	_, err := kubeConfigFn(thread, nil, nil, args)
	if err != nil {
		return err
	}

	return nil
}

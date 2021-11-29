// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// instanceListProvider is a built-in sarlark function that collects compute resources ad list of instance IDs.
// Starlark format: instance_list_provider(instances=<instance-list>, region="<aws-region>")
func instanceListProvider(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var instances *starlark.List
	var region string

	if err := starlark.UnpackArgs(
		identifiers.instanceListProvider, args, kwargs,
		"instances", &instances,
		"region", &region,
		); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.instanceListProvider, err)
	}

	if instances == nil || instances.Len() == 0 {
		return starlark.None, fmt.Errorf("%s: missing argument: instances", identifiers.instanceListProvider)
	}

	if region == "" {
		return starlark.None, fmt.Errorf("%s: missing argument region", identifiers.instanceListProvider)
	}

	cfgStruct := starlark.StringDict{
		"kind": starlark.String(identifiers.instanceListProvider),
		"transport": starlark.String("ssm"),
		"instances": instances,
		"region": starlark.String(region),
	}

	return starlarkstruct.FromStringDict(starlark.String(identifiers.instanceListProvider), cfgStruct), nil
}
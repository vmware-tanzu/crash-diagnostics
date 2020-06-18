// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// hostListProvider is a built-in starlark function that collects compute resources as a list of host IPs
// Starlark format: host_list_provider(hosts=<host-list> [, ssh_config=ssh_config()])
func hostListProvider(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var dictionary starlark.StringDict
	if kwargs != nil {
		dict, err := kwargsToStringDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		dictionary = dict
	}

	return newHostListProvider(thread, dictionary)
}

// newHostListProvider returns a struct with host list provider info
func newHostListProvider(thread *starlark.Thread, dictionary starlark.StringDict) (*starlarkstruct.Struct, error) {
	// validate args
	if _, ok := dictionary["hosts"]; !ok {
		return nil, fmt.Errorf("%s: missing hosts argument", identifiers.hostListProvider)
	}

	// augment args
	dictionary["kind"] = starlark.String(identifiers.hostListProvider)
	dictionary["transport"] = starlark.String("ssh")
	if _, ok := dictionary[identifiers.sshCfg]; !ok {
		data := thread.Local(identifiers.sshCfg)
		sshcfg, ok := data.(starlark.StringDict)
		if !ok {
			return nil, fmt.Errorf("%s: default ssh_config not found", identifiers.hostListProvider)
		}
		dictionary[identifiers.sshCfg] = starlarkstruct.FromStringDict(starlarkstruct.Default, sshcfg)
	}

	return starlarkstruct.FromStringDict(starlarkstruct.Default, dictionary), nil
}
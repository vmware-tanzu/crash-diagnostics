// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// resourcesFunc is a built-in starlark function that prepares returns compute list of resources.
// Starlark format: resources(provider=<provider-function>)
func resourcesFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if kwargs == nil {
		return starlark.None, fmt.Errorf("%s: missing arguments", identifiers.resources)
	}
	var dictionary starlark.StringDict
	if kwargs != nil {
		dict, err := kwargsToStringDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		dictionary = dict
	}

	var provider *starlarkstruct.Struct
	if hosts, ok := dictionary["hosts"]; ok {
		prov, err := newHostListProvider(thread, starlark.StringDict{"hosts": hosts})
		if err != nil {
			return starlark.None, err
		}
		provider = prov
	} else if prov, ok := dictionary["provider"]; ok {
		prov, ok := prov.(*starlarkstruct.Struct)
		if !ok {
			return starlark.None, fmt.Errorf("%s: provider not a struct", identifiers.resources)
		}
		provider = prov
	}

	if provider == nil {
		return starlark.None, fmt.Errorf("%s: hosts or provider argument required", identifiers.resources)
	}

	// enumerate resources from provider
	resources, err := enum(provider)
	if err != nil {
		return starlark.None, err
	}

	// save resources for future use
	thread.SetLocal(identifiers.resources, resources)

	return resources, nil
}

// enum returns a list of structs containing the fully enumerated compute resource
// info needed to execute commands.
func enum(provider *starlarkstruct.Struct) (*starlark.List, error) {
	if provider == nil {
		fmt.Errorf("missing provider")
	}

	var resources []starlark.Value

	kindVal, err := provider.Attr("kind")
	if err != nil {
		return nil, fmt.Errorf("provider missing field kind")
	}

	kind := trimQuotes(kindVal.String())

	switch kind {
	case identifiers.hostListProvider:
		hosts, err := provider.Attr("hosts")
		if err != nil {
			return nil, fmt.Errorf("hosts not found in %s", identifiers.hostListProvider)
		}

		hostList, ok := hosts.(*starlark.List)
		if !ok {
			return nil, fmt.Errorf("%s: unexpected type for hosts: %T", identifiers.hostListProvider, hosts)
		}

		transport, err := provider.Attr("transport")
		if err != nil {
			return nil, fmt.Errorf("transport not found in %s", identifiers.hostListProvider)
		}

		sshCfg, err := provider.Attr(identifiers.sshCfg)
		if err != nil {
			return nil, fmt.Errorf("ssh_config not found in %s", identifiers.hostListProvider)
		}

		for i := 0; i < hostList.Len(); i++ {
			dict := starlark.StringDict{
				"kind":       starlark.String(identifiers.hostResource),
				"provider":   starlark.String(identifiers.hostListProvider),
				"host":       hostList.Index(i),
				"transport":  transport,
				"ssh_config": sshCfg,
			}
			resources = append(resources, starlarkstruct.FromStringDict(starlarkstruct.Default, dict))
		}
	}

	return starlark.NewList(resources), nil
}

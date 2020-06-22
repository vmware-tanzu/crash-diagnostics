// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// addDefaultSshConf initalizes a Starlark Dict with default
// ssh_config configuration data
func addDefaultSSHConf(thread *starlark.Thread) error {
	args := makeDefaultSSHConfig()
	_, err := sshConfigFn(thread, nil, nil, args)
	if err != nil {
		return err
	}
	return nil
}

// sshConfigFn is the backing built-in fn that saves and returns its argument as struct value.
// Starlark format: ssh_config(conf0=val0, ..., confN=valN)
func sshConfigFn(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var dictionary starlark.StringDict
	if kwargs != nil {
		dict, err := kwargsToStringDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		dictionary = dict
	}

	// validation
	if _, ok := dictionary[identifiers.username]; !ok {
		return starlark.None, fmt.Errorf("%s: username required", identifiers.sshCfg)
	}
	if _, ok := dictionary[identifiers.port]; !ok {
		dictionary[identifiers.port] = starlark.String(defaults.sshPort)
	}
	if _, ok := dictionary[identifiers.maxRetries]; !ok {
		dictionary[identifiers.maxRetries] = starlark.MakeInt(defaults.connRetries)
	}
	if _, ok := dictionary[identifiers.privateKeyPath]; !ok {
		dictionary[identifiers.privateKeyPath] = starlark.String(defaults.pkPath)
	}

	structVal := starlarkstruct.FromStringDict(starlarkstruct.Default, dictionary)

	// save to be used as default when needed
	thread.SetLocal(identifiers.sshCfg, structVal)

	return structVal, nil
}

func makeDefaultSSHConfig() []starlark.Tuple {
	return []starlark.Tuple{
		starlark.Tuple{starlark.String("username"), starlark.String(getUsername())},
		starlark.Tuple{starlark.String("port"), starlark.String("22")},
		starlark.Tuple{starlark.String("private_key_path"), starlark.String(defaults.pkPath)},
		starlark.Tuple{starlark.String("max_retries"), starlark.MakeInt(defaults.connRetries)},
		starlark.Tuple{starlark.String("conn_timeout"), starlark.MakeInt(defaults.connTimeout)},
	}
}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"go.starlark.net/starlark"
)

// addDefaultSshConf initalizes a Starlark Dict with default
// ssh_config configuration data
func addDefaultSSHConf(thread *starlark.Thread) error {
	args := []starlark.Tuple{
		starlark.Tuple{starlark.String("username"), starlark.String(getUsername())},
		starlark.Tuple{starlark.String("private_key_path"), starlark.String(defaults.pkPath)},
		starlark.Tuple{starlark.String("conn_retries"), starlark.MakeInt(defaults.connRetries)},
		starlark.Tuple{starlark.String("conn_timeout"), starlark.MakeInt(defaults.connTimeout)},
	}

	_, err := sshConfigFn(thread, nil, nil, args)
	if err != nil {
		return err
	}

	return nil
}

// sshConfigFn is the backing built-in function for the `ssh_config` configuration function.
// It creates and returns a dictionary from collected configs (as kwargs)
// It also saves the dict into the thread as the last known ssh config to be used as default.
func sshConfigFn(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var dictionary *starlark.Dict
	if kwargs != nil {
		dict, err := tupleSliceToDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		dictionary = dict
	}

	// save to be used as default when needed
	thread.SetLocal(identifiers.sshCfg, dictionary)

	return dictionary, nil
}

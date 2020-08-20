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
// Starlark format: ssh_config(username=name[, port][, private_key_path][,max_retries][,conn_timeout][,jump_user][,jump_host])
func sshConfigFn(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var uname, port, pkPath, jUser, jHost string
	var maxRetries, connTimeout int

	if err := starlark.UnpackArgs(
		identifiers.crashdCfg, args, kwargs,
		"username", &uname,
		"port?", &port,
		"private_key_path?", &pkPath,
		"jump_user?", &jUser,
		"jump_host?", &jHost,
		"max_retries?", &maxRetries,
		"conn_timeout?", &connTimeout,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.hostListProvider, err)
	}

	// validation
	if len(uname) == 0 {
		return starlark.None, fmt.Errorf("%s: username required", identifiers.sshCfg)
	}
	if len(port) == 0 {
		port = defaults.sshPort
	}
	if maxRetries == 0 {
		maxRetries = defaults.connRetries
	}
	if connTimeout == 0 {
		connTimeout = defaults.connTimeout
	}
	if len(pkPath) == 0 {
		pkPath = defaults.pkPath
	}

	sshConfigDict := starlark.StringDict{
		"username":         starlark.String(uname),
		"port":             starlark.String(port),
		"private_key_path": starlark.String(pkPath),
		"max_retries":      starlark.MakeInt(maxRetries),
		"conn_timeout":     starlark.MakeInt(connTimeout),
	}
	if len(jUser) != 0 {
		sshConfigDict["jump_user"] = starlark.String(jUser)
	}
	if len(jHost) != 0 {
		sshConfigDict["jump_host"] = starlark.String(jHost)
	}
	structVal := starlarkstruct.FromStringDict(starlark.String(identifiers.sshCfg), sshConfigDict)

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

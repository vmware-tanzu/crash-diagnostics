// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sshconf

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("ssh_config")
	Func    = sshConfigFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	builtins.Register(Name, Builtin)
}

// sshConfigFunc implements a starlark built-in function that gathers ssh connection configuration.
//
// Example:
//    ssh_config(username="testuser", port="44")
//
// Args:
// - username (required)
// - port
// - private_key_path
// - jump_user
// - jump_host
// - max_retries
// - conn_timeout
func sshConfigFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	config := newCmd().Run(thread, args)

	// convert and return result
	return functions.Result(Name, config)
}

func DefaultSSHConfig() *Config {
	return &Config{
		Username:       functions.DefaultUsername(),
		Port:           DefaultPort(),
		PrivateKeyPath: DefaultPKPath(),
		MaxRetries:     DefaultMaxRetries(),
		ConnTimeout:    DefaultConnTimeout(),
	}
}

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package sshconf_alias implements an alias for the built-in function make_sshconf.Func.
package sshconf_alias

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf/make_sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	// Func is the built-in function that implements an alias for the make_sshconf.Func.
	// This alias returns the sshconf.Result.Config value directly as a convenience.
	// However, it will stop the script if an error occurs. For better error-handling
	// use the make_ssh_config function directly in scripts (implemented by make_sshconfig).
	Func       = sshConfigFunc
	Name       = functions.FunctionName("ssh_config")
	Builtin    = starlark.NewBuiltin(string(Name), Func)
	Identifier = string(Name)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// sshConfigFunc is the built-in function that implements an alias for the make_sshconf.Func.
// This alias returns the sshconf.Result.Config value directly as a convenience.
// However, it will stop the script if an error occurs. For better error-handling
// use the make_ssh_config function directly in scripts (implemented by make_sshconfig).
//
// Example:
//   ssh_config(username="testuser", port="44")
//
// Args:
// - username (required)
// - port
// - private_key_path
// - jump_user
// - jump_host
// - max_retries
// - conn_timeout
//
// Returns
// - Config: a sshconf.Config containing the configuration data
//
func sshConfigFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args sshconf.Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := make_sshconf.Run(thread, args)

	// return fatal error/stop script
	if result.Error != "" {
		return starlark.None, errors.New(result.Error)
	}

	// convert and return result
	return functions.Result(Name, result.Config)
}

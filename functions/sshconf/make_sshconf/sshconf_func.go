// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package make_sshconf

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("make_ssh_config")
	Func    = sshConfigFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
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
	var args sshconf.Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	if args.Username == "" {
		functions.Error(Name, fmt.Errorf("%s: username is empty", Name))
	}

	result := Run(thread, args)

	// convert and return result
	return functions.Result(Name, result)
}

// Run executes command function
func Run(t *starlark.Thread, args sshconf.Args) sshconf.Result {
	if args.Port == "" {
		args.Port = sshconf.DefaultPort()
	}
	if args.PrivateKeyPath == "" {
		args.PrivateKeyPath = sshconf.DefaultPKPath()
	}
	if args.ConnTimeout == 0 {
		args.ConnTimeout = sshconf.DefaultConnTimeout()
	}

	// add private key to agent if agent was saved in thread
	if agent, ok := sshconf.SSHAgentFromThread(t); ok {
		if err := agent.AddKey(args.PrivateKeyPath); err != nil {
			return sshconf.Result{Error: fmt.Sprintf("unable to add private key to agent: %s", args.PrivateKeyPath)}
		}
	}

	return sshconf.Result{
		Config: sshconf.Config{
			Username:       args.Username,
			Port:           args.Port,
			PrivateKeyPath: args.PrivateKeyPath,
			JumpUsername:   args.JumpUsername,
			JumpHost:       args.JumpHost,
			MaxRetries:     args.MaxRetries,
			ConnTimeout:    args.ConnTimeout,
		},
	}
}

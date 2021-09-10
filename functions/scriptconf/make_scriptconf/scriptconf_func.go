// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package make_scriptconf

import (
	"fmt"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"github.com/vmware-tanzu/crash-diagnostics/util"
	"go.starlark.net/starlark"
)

var (
	Name       = functions.FunctionName("make_script_config")
	Func       = makeScriptConfigFunc
	Builtin    = starlark.NewBuiltin(string(Name), Func)
	Identifier = string(Name)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// makeScriptConfigFunc implements a starlark built-in function that gathers and stores configuration
// settings for a running script. This function returns function Result that contains either an error,
// if one occured, or a Config value.
//
// Example:
//    make_script_config(workdir=path, default_shell=shellpath, requires=["command0",...,"commandN"])
//
// Args:
//   - workdir string - a path that can be used as work directory during script exec
//   - gid string - the default group id to use when executing an OS command
//   - uid string - a default userid to use when executing an OS command
//   - default_shell string - path to a shell program that can be used as default (i.e. /bin/sh)
//   - requires [] string - a list of paths for commands that should be on the machine where script is executed
//   - use_ssh_agent bool - specifies if an ssh-agent should be setup for private key management
//
// Returns
//  - Error: an error message if the call generated one
//  - Config: a scriptconf.Config containing the configuration data
func makeScriptConfigFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args scriptconf.Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := Run(thread, args)

	// save/overwrite config result in thread
	thread.SetLocal(scriptconf.Identifier, result.Config)

	// convert and return result
	return functions.Result(Name, result)
}

// Run executes the command function
func Run(t *starlark.Thread, args scriptconf.Args) scriptconf.Result {
	if err := validateArgs(&args); err != nil {
		return scriptconf.Result{Error: fmt.Sprintf("failed to validate configuration: %s", err)}
	}

	// create workdir if needed
	if err := functions.MakeDir(args.Workdir, 0744); err != nil && !os.IsExist(err) {
		return scriptconf.Result{Error: fmt.Sprintf("failed to create workdir: %s", err)}
	}

	// start local ssh-agent
	if args.UseSSHAgent {
		_, err := sshconf.MakeDefaultSSHAgentForThread(t)
		if err != nil {
			return scriptconf.Result{Error: fmt.Sprintf("%s: failed to start ssh agent: %s", string(Name), err)}
		}
	}

	return scriptconf.Result{
		Config: scriptconf.Config{
			Workdir:      args.Workdir,
			Gid:          args.Gid,
			Uid:          args.Uid,
			DefaultShell: args.DefaultShell,
			Requires:     args.Requires,
			UseSSHAgent:  args.UseSSHAgent,
		},
	}
}

func validateArgs(params *scriptconf.Args) error {
	if params.Workdir == "" {
		params.Workdir = scriptconf.DefaultWorkdir()
	}
	wd, err := util.ExpandPath(params.Workdir)
	if err != nil {
		return err
	}
	params.Workdir = wd

	if params.Gid == "" {
		params.Gid = functions.DefaultGid()
	}

	if params.Uid == "" {
		params.Uid = functions.DefaultUid()
	}

	return nil
}

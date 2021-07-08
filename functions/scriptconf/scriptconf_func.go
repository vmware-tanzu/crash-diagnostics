// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"errors"
	"fmt"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"github.com/vmware-tanzu/crash-diagnostics/util"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("script_config")
	Func    = scriptConfigFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// scriptConfigFunc implements a starlark built-in function that gathers and stores configuration
// settings for a running script.
//
// Example:
//    script_config(workdir=path, default_shell=shellpath, requires=["command0",...,"commandN"])
//
// Args:
//   - workdir string - a path that can be used as work directory during script exec
//   - gid string - the default group id to use when executing an OS command
//   - uid string - a default userid to use when executing an OS command
//   - default_shell string - path to a shell program that can be used as default (i.e. /bin/sh)
//   - requires [] string - a list of paths for commands that should be on the machine where script is executed
//   - use_ssh_agent bool - specifies if an ssh-agent should be setup for private key management
func scriptConfigFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := Run(thread, args)

	// save config result in thread
	thread.SetLocal(string(Name), result)

	// convert and return result
	return functions.Result(Name, result)
}

// Run executes the command function
func Run(t *starlark.Thread, args Args) Result {
	if err := validateArgs(&args); err != nil {
		return Result{Error: fmt.Sprintf("failed to validate configuration: %s", err)}
	}

	// create workdir if needed
	if err := functions.MakeDir(args.Workdir, 0744); err != nil && !os.IsExist(err) {
		return Result{Error: fmt.Sprintf("failed to create workdir: %s", err)}
	}

	// start local ssh-agent
	if args.UseSSHAgent {
		_, err := sshconf.MakeSSHAgentForThread(t)
		if err != nil {
			return Result{Error: fmt.Sprintf("%s: failed to start ssh agent: %s", string(Name), err)}
		}
	}

	return Result{
		Config: Config{
			Workdir:      args.Workdir,
			Gid:          args.Gid,
			Uid:          args.Uid,
			DefaultShell: args.DefaultShell,
			Requires:     args.Requires,
			UseSSHAgent:  args.UseSSHAgent,
		},
	}
}

// ConfigFromThread retrieves script config result from provided
// thread instance. If found, bool = true.
func ConfigFromThread(t *starlark.Thread) (Config, bool) {
	if val := t.Local(string(Name)); val != nil {
		result, ok := val.(Config)
		if !ok {
			return Config{}, ok
		}
		return result, true
	}
	return Config{}, false
}

func MakeConfigForThread(t *starlark.Thread) (Config, error) {
	conf := makeDefaultConf()
	args := Args{
		Workdir:      conf.Workdir,
		Gid:          conf.Gid,
		Uid:          conf.Uid,
		DefaultShell: conf.DefaultShell,
		Requires:     conf.Requires,
		UseSSHAgent:  conf.UseSSHAgent,
	}
	result := Run(t, args)
	if result.Error != "" {
		return Config{}, errors.New(result.Error)
	}
	return result.Config, nil
}

func makeDefaultConf() Config {
	return Config{
		Workdir:      DefaultWorkdir(),
		Gid:          functions.DefaultGid(),
		Uid:          functions.DefaultUid(),
		DefaultShell: "/bin/sh",
		Requires:     []string{"/bin/ssh", "/bin/scp"},
		UseSSHAgent:  false,
	}
}

func validateArgs(params *Args) error {
	if params.Workdir == "" {
		params.Workdir = DefaultWorkdir()
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

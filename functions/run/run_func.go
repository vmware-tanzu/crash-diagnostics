// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("run")
	Func    = runFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	builtins.Register(Name, Builtin)
}

// runFunc implements a starlark built-in function `run()` that can execute processes on remote
// compute resource.
//
// Example:
//    run(cmd="echo 'hello'", resources=hostlist_provider(hosts=["host1","host2"]))
//
// Args:
// - cmd: the command to run (required)
// - ssh_config: ssh configuration
// - resources: compute resources to run command
func runFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	if args.Cmd == "" {
		return functions.Error(Name, fmt.Errorf("%s: missing command", Name))
	}

	if args.Resources == nil {
		res, ok := providers.ResourcesFromThread(thread)
		if !ok {
			return functions.Error(Name, fmt.Errorf("%s: missing resources", Name))
		}
		args.Resources = &res
	}

	if args.SSHConfig == nil {
		args.SSHConfig = sshconf.DefaultSSHConfig()
	}

	agent, ok := sshconf.SSHAgentFromThread(thread)
	if !ok {
		return functions.Error(Name, fmt.Errorf("%s: missing ssh-agent instance", Name))
	}

	result := newCmd().Run(thread, agent, args)

	// convert and return result
	return functions.Result(Name, result)
}

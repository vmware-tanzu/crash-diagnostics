// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
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
// - resources: list of compute resources to run command
func runFunc(thread *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	if args.Cmd == "" {
		return functions.Error(Name, fmt.Errorf("%s: missing command", Name))
	}

	if reflect.ValueOf(args.Resources).IsZero() {
		res, ok := providers.ResourcesFromThread(thread)
		if !ok {
			return functions.Error(Name, fmt.Errorf("%s: missing resources", Name))
		}
		args.Resources = res
	}

	if reflect.ValueOf(args.SSHConfig).IsZero() {
		conf := sshconf.DefaultConfig()
		args.SSHConfig = conf
	}

	// check for ssh-agent
	agent, ok := sshconf.SSHAgentFromThread(thread)
	if !ok {
		// is there a script config
		conf, scOk := scriptconf.ConfigFromThread(thread)
		if scOk && conf.UseSSHAgent { // no script config, bail
			return functions.Error(Name, fmt.Errorf("%s: ssh-agent not found", Name))
		} else {
			logrus.Warnf("%s: not using ssh-agent", Name)
		}
	}

	result := Run(thread, agent, args)

	// convert and return result
	return functions.Result(Name, result)
}

// Run runs the command function
func Run(t *starlark.Thread, agent ssh.Agent, args Args) Result {
	if reflect.ValueOf(args.SSHConfig).IsZero() {
		return Result{Error: "missing SSH config"}
	}
	sshConf := args.SSHConfig

	if reflect.ValueOf(args.Resources).IsZero() {
		return Result{Error: "missing resources"}
	}
	hosts := args.Resources.Hosts
	if len(hosts) == 0 {
		return Result{Error: fmt.Sprintf("%s provided no host", args.Resources.Provider)}
	}

	var cmdResults []RemoteProc
	for _, host := range hosts {
		var jumpProxy *ssh.ProxyJumpArgs
		if sshConf.JumpHost != "" && sshConf.JumpUsername != "" {
			jumpProxy = &ssh.ProxyJumpArgs{
				User: sshConf.JumpUsername,
				Host: sshConf.JumpHost,
			}
		}
		sshArgs := ssh.SSHArgs{
			User:           sshConf.Username,
			Host:           host,
			Port:           sshConf.Port,
			MaxRetries:     int(sshConf.MaxRetries),
			ProxyJump:      jumpProxy,
			PrivateKeyPath: sshConf.PrivateKeyPath,
		}

		var errMsg string
		sshOutput, err := ssh.Run(sshArgs, agent, args.Cmd)
		if err != nil {
			errMsg = err.Error()
		}
		cmdResults = append(cmdResults, RemoteProc{Error: errMsg, Host: host, Output: sshOutput})
	}
	return Result{Procs: cmdResults}
}

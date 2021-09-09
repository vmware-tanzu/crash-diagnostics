// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package copy_to

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("copy_to")
	Func    = copyToFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register Starlark built-in
func init() {
	registrar.Register(Name, Builtin)
}

// copyToFunc is a built-in starlark function that copies specified resources from the local machine to a remote compute resource.
// Starlark format: result = copy_to(source_path=<file path>, target_path=<file path>, ssh_config=<config>, resources=resources])
//
// Args:
// - source_path: local source path of file to copy
// - target_path: remote target path of file to copy (option)
// - ssh_config: ssh configuration
// - resources: list of compute resources from which to copy
//
func copyToFunc(thread *starlark.Thread, b *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	if args.SourcePath == "" {
		return functions.Error(Name, fmt.Errorf("%s: missing source path", Name))
	}

	if args.TargetPath == "" {
		args.TargetPath = args.SourcePath
	}

	// retrieve resources from thread if none is provided
	if reflect.ValueOf(args.Resources).IsZero() {
		res, ok := providers.ResourcesFromThread(thread)
		if !ok {
			return functions.Error(Name, fmt.Errorf("%s: missing resources", Name))
		}
		args.Resources = res
	}

	if reflect.ValueOf(args.SSHConfig).IsZero() || reflect.ValueOf(args.SSHConfig).IsZero() {
		// attempt to get it from thread, else return default
		conf, ok := sshconf.ConfigFromThread(thread)
		if !ok || reflect.ValueOf(conf).IsZero() {
			conf = sshconf.DefaultConfig()
		}
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

func Run(_ *starlark.Thread, agent ssh.Agent, args Args) Result {
	sshConf := args.SSHConfig
	hosts := args.Resources.Hosts
	if len(hosts) == 0 {
		return Result{Error: fmt.Sprintf("%s provided no host", args.Resources.Provider)}
	}

	var jumpProxy *ssh.ProxyJumpArgs
	if sshConf.JumpHost != "" && sshConf.JumpUsername != "" {
		jumpProxy = &ssh.ProxyJumpArgs{
			User: sshConf.JumpUsername,
			Host: sshConf.JumpHost,
		}
	}

	var copies []CopyResult
	for _, host := range hosts {

		sshArgs := ssh.SSHArgs{
			User:           sshConf.Username,
			Host:           host,
			Port:           sshConf.Port,
			MaxRetries:     int(sshConf.MaxRetries),
			ProxyJump:      jumpProxy,
			PrivateKeyPath: sshConf.PrivateKeyPath,
		}

		var errMsg string
		err := ssh.CopyTo(sshArgs, agent, args.SourcePath, args.TargetPath)
		if err != nil {
			errMsg = err.Error()
		}

		copies = append(copies, CopyResult{Error: errMsg, Host: host, SourcePath: args.SourcePath, TargetPath: args.TargetPath})
	}

	return Result{Copies: copies}
}

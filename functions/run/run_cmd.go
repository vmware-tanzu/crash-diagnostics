// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"go.starlark.net/starlark"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

func (c *cmd) Run(t *starlark.Thread, agent ssh.Agent, args Args) Result {
	if agent == nil {
		return Result{Error: "missing ssh-agent"}
	}
	sshConf := args.SSHConfig
	if sshConf == nil {
		return Result{Error: "missing SSH config"}
	}

	if args.Resources == nil {
		return Result{Error: "missing resources"}
	}
	hosts := args.Resources.Hosts
	if len(hosts) == 0 {
		return Result{Error: fmt.Sprintf("%s provided no host", args.Resources.Provider)}
	}

	var cmdResults []CmdResult
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
		cmdResults = append(cmdResults, CmdResult{Error: errMsg, Host: host, Output: sshOutput})
	}
	return Result{CmdResults: cmdResults}
}

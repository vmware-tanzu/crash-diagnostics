// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"go.starlark.net/starlark"
)

func TestCmd_Run(t *testing.T) {
	sshAgent, err := ssh.StartAgent()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		setup func(*testing.T) Args
		eval  func(*testing.T, Args)
	}{
		{
			name: "missing resources",
			setup: func(t *testing.T) Args {
				return Args{Cmd: "echo 'Hello World!'", SSHConfig: sshconf.DefaultSSHConfig()}
			},
			eval: func(t *testing.T, args Args) {
				result := newCmd().Run(&starlark.Thread{}, sshAgent, args)
				if result.Error == "" {
					t.Error("expecting error, got none")
				}
			},
		},
		{
			name: "missing SSHConfig",
			setup: func(t *testing.T) Args {
				return Args{Cmd: "echo 'Hello World!'", Resources: &providers.Resources{}}
			},
			eval: func(t *testing.T, args Args) {
				result := newCmd().Run(&starlark.Thread{}, sshAgent, args)
				if result.Error == "" {
					t.Error("expecting error, got none")
				}
			},
		},
		{
			name: "missing ssh-agent",
			setup: func(t *testing.T) Args {
				return Args{Cmd: "echo 'Hello World!'", SSHConfig: sshconf.DefaultSSHConfig(), Resources: &providers.Resources{}}
			},
			eval: func(t *testing.T, args Args) {
				result := newCmd().Run(&starlark.Thread{}, nil, args)
				if result.Error == "" {
					t.Error("expecting error, got none")
				}
			},
		},
		{
			name: "simple cmd",
			setup: func(t *testing.T) Args {
				return Args{
					Cmd:       "echo 'Hello World!'",
					Resources: &providers.Resources{Hosts: []string{"127.0.0.1"}},
					SSHConfig: sshconf.DefaultSSHConfig(),
				}
			},
			eval: func(t *testing.T, args Args) {
				args.SSHConfig.Username = testSupport.CurrentUsername()
				args.SSHConfig.Port = testSupport.PortValue()
				args.SSHConfig.PrivateKeyPath = testSupport.PrivateKeyPath()
				args.SSHConfig.MaxRetries = int64(testSupport.MaxConnectionRetries())

				result := newCmd().Run(&starlark.Thread{}, sshAgent, args)
				if result.Error != "" {
					t.Error(result.Error)
				}
				if len(result.CmdResults) != 1 {
					t.Error("missing command result")
				}
				output := strings.TrimSpace(result.CmdResults[0].Output)
				if output != "Hello World!" {
					t.Error("unexpected result:", output)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.setup(t))
		})
	}

	if err := sshAgent.Stop(); err != nil {
		t.Fatal(err)
	}
}

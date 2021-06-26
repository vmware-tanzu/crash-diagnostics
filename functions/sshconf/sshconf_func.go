// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sshconf

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
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

	if args.Username == "" {
		functions.Error(Name, fmt.Errorf("%s: username is empty", Name))
	}

	config := Run(thread, args)

	// convert and return result
	return functions.Result(Name, config)
}

// Run executes command function
func Run(t *starlark.Thread, args Args) Result {
	if args.Port == "" {
		args.Port = DefaultPort()
	}
	if args.PrivateKeyPath == "" {
		args.PrivateKeyPath = DefaultPKPath()
	}
	if args.ConnTimeout == 0 {
		args.ConnTimeout = DefaultConnTimeout()
	}

	// add private key to agent if agent was saved in thread
	if agent, ok := SSHAgentFromThread(t); ok {
		if err := agent.AddKey(args.PrivateKeyPath); err != nil {
			return Result{Error: fmt.Sprintf("unable to add private key to agent: %s", args.PrivateKeyPath)}
		}
	}

	return Result{
		Conf: Config{
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

func SSHAgentFromThread(t *starlark.Thread) (ssh.Agent, bool) {
	if agentVal := t.Local(AgentIdentifier); agentVal != nil {
		agent, ok := agentVal.(ssh.Agent)
		if !ok {
			return nil, false
		}
		return agent, true
	}
	return nil, false
}

// ConfigFromThread returns an sshconf.Config from provided
// starlark thread.
func ConfigFromThread(t *starlark.Thread) (Config, bool) {
	if confVal := t.Local(Identifier); confVal != nil {
		conf, ok := confVal.(Config)
		if !ok {
			return Config{}, false
		}
		return conf, true
	}
	return Config{}, false
}

func MakeConfigForThread(t *starlark.Thread) (Config, error) {
	conf := makeDefaultSSHConfig()
	args := Args{
		Username:       conf.Username,
		Port:           conf.Port,
		PrivateKeyPath: conf.PrivateKeyPath,
		MaxRetries:     conf.MaxRetries,
		ConnTimeout:    conf.ConnTimeout,
	}
	result := Run(t, args)
	if result.Error != "" {
		return Config{}, errors.New(result.Error)
	}
	return result.Conf, nil
}

func MakeSSHAgentForThread(t *starlark.Thread) (ssh.Agent, error) {
	agent, err := ssh.StartAgent()
	if err != nil {
		return nil, err
	}
	t.SetLocal(AgentIdentifier, agent)
	return agent, nil
}

func makeDefaultSSHConfig() Config {
	return Config{
		Username:       functions.DefaultUsername(),
		Port:           DefaultPort(),
		PrivateKeyPath: DefaultPKPath(),
		MaxRetries:     DefaultMaxRetries(),
		ConnTimeout:    DefaultConnTimeout(),
	}
}

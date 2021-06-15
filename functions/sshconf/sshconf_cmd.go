// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sshconf

import (
	"fmt"

	"go.starlark.net/starlark"
)

type confCmd struct{}

func newCmd() *confCmd {
	return new(confCmd)
}

// Run applies processes the params and generates a configuration
func (c *confCmd) Run(t *starlark.Thread, args Args) Config {
	if args.Username == "" {
		return Config{Error: "username required"}
	}
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
			return Config{Error: fmt.Sprintf("unable to add private key to agent: %s", args.PrivateKeyPath)}
		}
	}

	return Config{
		Username:       args.Username,
		Port:           args.Port,
		PrivateKeyPath: args.PrivateKeyPath,
		JumpUsername:   args.JumpUsername,
		JumpHost:       args.JumpHost,
		MaxRetries:     args.MaxRetries,
		ConnTimeout:    args.ConnTimeout,
	}
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
)

// AuthConfigCommand represents AUTHCONFIG directive:
//
// AUTHCONFIG username:"username" private-key:"/path/to/key api-key:key"
//
// Param username is required
type AuthConfigCommand struct {
	cmd
}

// NewAuthConfigCommand parses the rawArgs and returns an  *AuthCommand
func NewAuthConfigCommand(index int, rawArgs string) (*AuthConfigCommand, error) {
	if err := validateRawArgs(CmdAuthConfig, rawArgs); err != nil {
		return nil, err
	}

	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("AUTHCONFIG: %v", err)
	}
	if err := validateCmdArgs(CmdAuthConfig, argMap); err != nil {
		return nil, err
	}

	cmd := &AuthConfigCommand{cmd: cmd{index: index, name: CmdAuthConfig, args: argMap}}

	return cmd, nil
}

// Index is the position of the command in the script
func (c *AuthConfigCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *AuthConfigCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *AuthConfigCommand) Args() map[string]string {
	return c.cmd.args
}

// GetPrivateKey returns the path of the private key configured
func (c *AuthConfigCommand) GetPrivateKey() string {
	return ExpandEnv(c.cmd.args["private-key"])
}

// GetUsername returns the User ID configured
func (c *AuthConfigCommand) GetUsername() string {
	return ExpandEnv(c.cmd.args["username"])
}

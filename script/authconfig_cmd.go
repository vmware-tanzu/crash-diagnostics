// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"strings"
)

// AuthConfigCommand captures auth configuration for a runtime
type AuthConfigCommand struct {
	cmd
	username   string
	privateKey string
	apiKey     string
}

// NewAuthConfigCommand parses the args and return a value of type *AuthCommand from:
// AUTHCONFIG username:<user-name> private-key:<path/to/key> api-key:<api-key-value>
func NewAuthConfigCommand(index int, args []string) (*AuthConfigCommand, error) {
	cmd := &AuthConfigCommand{cmd: cmd{index: index, name: CmdAuthConfig, args: args}}

	if err := validateCmdArgs(CmdAuthConfig, args); err != nil {
		return nil, err
	}

	// split each arg
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s: bad argument: %s", CmdAuthConfig, arg)
		}

		switch {
		case strings.EqualFold(parts[0], "username"):
			cmd.username = parts[1]
		case strings.EqualFold(parts[0], "private-key"):
			cmd.privateKey = parts[1]
		case strings.EqualFold(parts[0], "api-key"):
			cmd.apiKey = parts[1]
		default:
		}

	}

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
func (c *AuthConfigCommand) Args() []string {
	return c.cmd.args
}

// GetPrivateKey returns the path of the private key configured
func (c *AuthConfigCommand) GetPrivateKey() string {
	return c.privateKey
}

// GetUsername returns the User ID configured
func (c *AuthConfigCommand) GetUsername() string {
	return c.username
}

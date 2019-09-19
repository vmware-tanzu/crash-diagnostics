// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"strings"
)

// SSHConfigCommand represents an SSHCONFIG directive in a script
type SSHConfigCommand struct {
	cmd
	userid     string
	privateKey string
}

// NewSSHConfigCommand parses the args and return a value of type *SSHCommand
func NewSSHConfigCommand(index int, args []string) (*SSHConfigCommand, error) {
	cmd := &SSHConfigCommand{cmd: cmd{index: index, name: CmdSSHConfig, args: args}}

	if err := validateCmdArgs(CmdSSHConfig, args); err != nil {
		return nil, err
	}

	sshParts := strings.Split(args[0], ":")
	switch {
	case len(sshParts) > 1:
		cmd.userid = sshParts[0]
		cmd.privateKey = sshParts[1]
	case len(sshParts) == 1:
		cmd.privateKey = sshParts[0]
	default:
		return nil, fmt.Errorf("SSHCONFIG misconfigured, expects SSHCONFIG [userid:]/path/to/key")
	}
	return cmd, nil
}

// Index is the position of the command in the script
func (c *SSHConfigCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *SSHConfigCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *SSHConfigCommand) Args() []string {
	return c.cmd.args
}

// GetPrivateKeyPath returns the path of the private key configured
func (c *SSHConfigCommand) GetPrivateKeyPath() string {
	return c.privateKey
}

// GetUserId returns the User ID configured
func (c *SSHConfigCommand) GetUserId() string {
	return c.userid
}

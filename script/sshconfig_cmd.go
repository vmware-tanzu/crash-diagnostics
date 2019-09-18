// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"strings"
)

// SSHConfigCommand represents an SSH configuration
// used for remote execution of CLI commands.
type SSHConfigCommand struct {
	cmd
	userid     string
	privateKey string
}

// NewSSHConfigCommand creates a value of type SSHCommand
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

func (c *SSHConfigCommand) Index() int {
	return c.cmd.index
}

func (c *SSHConfigCommand) Name() string {
	return c.cmd.name
}

func (c *SSHConfigCommand) Args() []string {
	return c.cmd.args
}

func (c *SSHConfigCommand) GetPrivateKeyPath() string {
	return c.privateKey
}

func (c *SSHConfigCommand) GetUserId() string {
	return c.userid
}

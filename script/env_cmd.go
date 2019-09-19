// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	envSep = regexp.MustCompile(`=`)
)

// EnvCommand represents ENV in a script
type EnvCommand struct {
	cmd
	envs []string
}

// NewEnvCommand returns parses the args as environment variables and returns *EnvCommand
func NewEnvCommand(index int, args []string) (*EnvCommand, error) {
	cmd := &EnvCommand{cmd: cmd{index: index, name: CmdEnv, args: args}}

	if err := validateCmdArgs(CmdEnv, args); err != nil {
		return nil, err
	}

	for _, arg := range args {
		parts := envSep.Split(strings.TrimSpace(arg), -1)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid ENV arg %s", arg)
		}
		cmd.envs = append(cmd.envs, fmt.Sprintf("%s=%s", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])))
	}

	return cmd, nil
}

// Index is the position of the command in the script
func (c *EnvCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *EnvCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *EnvCommand) Args() []string {
	return c.cmd.args
}

// Envs returns slice of the parsed declared environment variables
func (c *EnvCommand) Envs() []string {
	return c.envs
}

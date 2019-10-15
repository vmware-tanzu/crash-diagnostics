// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	envSep = regexp.MustCompile(`=`)
)

// EnvCommand represents ENV directive:
//
// ENV [envs:"]key0=val0...keyN=valN["]
//
// Supports multiple ENV in one script.
type EnvCommand struct {
	cmd
	envs map[string]string
}

// NewEnvCommand returns parses the args as environment variables and returns *EnvCommand
func NewEnvCommand(index int, rawArgs string) (*EnvCommand, error) {
	if err := validateRawArgs(CmdEnv, rawArgs); err != nil {
		return nil, err
	}

	// by default the args are stored in envs
	argMap := map[string]string{"envs": rawArgs}
	cmd := &EnvCommand{
		envs: make(map[string]string),
		cmd:  cmd{index: index, name: CmdEnv, args: argMap},
	}

	if err := validateCmdArgs(CmdEnv, argMap); err != nil {
		return nil, err
	}

	envs := spaceSep.Split(rawArgs, -1)
	for _, env := range envs {
		parts := envSep.Split(strings.TrimSpace(env), -1)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid ENV arg %s", env)
		}
		key, val := parts[0], parts[1]
		cmd.envs[key] = val
		os.Setenv(parts[0], parts[1])
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
func (c *EnvCommand) Args() map[string]string {
	return c.cmd.args
}

// Envs returns slice of the parsed declared environment variables
func (c *EnvCommand) Envs() map[string]string {
	return c.envs
}

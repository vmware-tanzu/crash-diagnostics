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

// EnvCommand represents ENV directive which can have one of the following forms:
//
//     ENV var0=val0 var1=val0 ... varN=valN
//     ENV vars:"var0=val0 var1=val0 ... varN=valN"
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

	// map params
	var argMap map[string]string
	if !isNamedParam(rawArgs) {
		rawArgs = makeNamedPram("vars", rawArgs)
	}
	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("ENV: %v", err)
	}

	cmd := &EnvCommand{
		envs: make(map[string]string),
		cmd:  cmd{index: index, name: CmdEnv, args: argMap},
	}

	if err := validateCmdArgs(CmdEnv, argMap); err != nil {
		return nil, err
	}

	// supported format keyN=valN keyN="valN" keyN='valN'
	// foreach key0=val0 key1=val1 ... keyN=valN
	// split into keyN, valN
	envs, err := wordSplit(argMap["vars"])
	if err != nil {
		return nil, fmt.Errorf("ENV: %s", err)
	}

	for _, env := range envs {
		parts := envSep.Split(strings.TrimSpace(env), 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("ENV: invalid: %s", env)
		}

		key := parts[0]
		val, err := wordSplit(parts[1]) // safely remove outer quotes
		if err != nil {
			return nil, fmt.Errorf("ENV: %s", err)
		}
		value := val[0]

		cmd.envs[key] = ExpandEnv(value)
		if err := os.Setenv(key, ExpandEnv(value)); err != nil {
			return nil, fmt.Errorf("ENV: %s", err)
		}
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

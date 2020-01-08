// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
)

// RunCommand represents RUN directive which
// can have one of the following two forms as shown below:
//
//     RUN <command-string>
//     RUN cmd:"<command-string>" shell:"shell-path" desc:"cmd-desc"
//
// The former takes no named parameter. When the latter form is used,
// parameter cmd: is required.
type RunCommand struct {
	cmd
}

// NewRunCommand returns *RunCommand with parsed arguments
func NewRunCommand(index int, rawArgs string) (*RunCommand, error) {
	if err := validateRawArgs(CmdRun, rawArgs); err != nil {
		return nil, err
	}

	// determine args
	argMap := make(map[string]string)
	if !isNamedParam(rawArgs) {
		// setup default param
		if isQuoted(rawArgs) {
			argMap["cmd"] = trimQuotes(rawArgs)
		} else {
			argMap["cmd"] = rawArgs
		}
	} else {
		args, err := mapArgs(rawArgs)
		if err != nil {
			return nil, fmt.Errorf("RUN: %v", err)
		}
		argMap = args
	}

	if err := validateCmdArgs(CmdRun, argMap); err != nil {
		return nil, fmt.Errorf("RUN: %s", err)
	}

	cmd := &RunCommand{cmd: cmd{index: index, name: CmdCapture, args: argMap}}
	return cmd, nil
}

// Index is the position of the command in the script
func (c *RunCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *RunCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *RunCommand) Args() map[string]string {
	return c.cmd.args
}

// GetCmdShell returns shell program and arguments
// for running the command string (i.e. /bin/bash -c)
func (c *RunCommand) GetCmdShell() string {
	return ExpandEnv(c.cmd.args["shell"])
}

// GetCmdString returns the raw CLI command string
func (c *RunCommand) GetCmdString() string {
	return ExpandEnv(c.cmd.args["cmd"])
}

// GetEffectiveCmd returns the shell (if any) and command as
// a slice of strings
func (c *RunCommand) GetEffectiveCmd() ([]string, error) {
	cmdStr := c.GetCmdString()
	shell := c.GetCmdShell()
	if c.GetCmdShell() != "" {
		shArgs, err := wordSplit(shell)
		if err != nil {
			return nil, err
		}
		return append(shArgs, cmdStr), nil
	}
	cmdArgs, err := wordSplit(cmdStr)
	if err != nil {
		return nil, err
	}
	return cmdArgs, nil
}

// GetParsedCmd returns the effective parsed command as commandName
// followed by a slice of command arguments
func (c *RunCommand) GetParsedCmd() (string, []string, error) {
	args, err := c.GetEffectiveCmd()
	if err != nil {
		return "", nil, err
	}
	return args[0], args[1:], nil
}

// GetEffectiveCmdStr returns the effective command as a string
// which wraps the command around a shell quote if necessary
func (c *RunCommand) GetEffectiveCmdStr() (string, error) {
	cmdStr := c.GetCmdString()
	shell := c.GetCmdShell()
	if c.GetCmdShell() != "" {
		return fmt.Sprintf("%s %s", shell, quote(cmdStr)), nil
	}
	return cmdStr, nil
}

// GetEcho returns the echo param for command. When
// set to {yes|true|on} the result of the command will be
// redirected to the stdout|stderr
func (c *RunCommand) GetEcho() string {
	return ExpandEnv(c.cmd.args["echo"])
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
)

// CaptureCommand represents CAPTURE directive which
// can have one of the following two forms as shown below:
//
//     CAPTURE <command-string>
//     CAPTURE cmd:"<command-string>" name:"cmd-name" desc:"cmd-desc"
//
// The former takes no named parameter. When the latter form is used,
// parameter cmd: is required.
type CaptureCommand struct {
	cmd
	cmdName string
	cmdArgs []string
}

// NewCaptureCommand returns *CaptureCommand with parsed arguments
func NewCaptureCommand(index int, rawArgs string) (*CaptureCommand, error) {
	if err := validateRawArgs(CmdCapture, rawArgs); err != nil {
		return nil, err
	}

	// determine args
	var argMap map[string]string
	if !isNamedParam(rawArgs) {
		// setup default param (notice quoted value)
		rawArgs = makeNamedPram("cmd", rawArgs)
	}
	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("CAPTURE: %v", err)
	}

	if err := validateCmdArgs(CmdCapture, argMap); err != nil {
		return nil, fmt.Errorf("CAPTURE: %s", err)
	}

	cmd := &CaptureCommand{cmd: cmd{index: index, name: CmdCapture, args: argMap}}

	cmdName, cmdArgs, err := cmdParse(cmd.GetCmdString())
	if err != nil {
		return nil, fmt.Errorf("CAPTURE: %s", err)
	}
	cmd.cmdName = cmdName
	cmd.cmdArgs = cmdArgs
	return cmd, nil
}

// Index is the position of the command in the script
func (c *CaptureCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *CaptureCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *CaptureCommand) Args() map[string]string {
	return c.cmd.args
}

// GetCmdString returns the raw CLI command string
func (c *CaptureCommand) GetCmdString() string {
	return c.cmd.args["cmd"]
}

// GetParsedCmd returns the CLI command name to be captured and its arguments
func (c *CaptureCommand) GetParsedCmd() (string, []string) {
	return c.cmdName, c.cmdArgs
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import "fmt"

// CaptuerCommand represents CAPTURE directive:
//
// CAPTURE cmd:<cmd-string> [shell:/path/to/shell]
//
type CaptureCommand struct {
	cmd
}

// NewCaptureCommand returns *CaptureCommand with parsed arguments
func NewCaptureCommand(index int, rawArgs string) (*CaptureCommand, error) {
	if err := validateRawArgs(CmdCapture, rawArgs); err != nil {
		return nil, err
	}

	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("CAPTURE: %v", err)
	}

	cmd := &CaptureCommand{cmd: cmd{index: index, name: CmdCapture, args: argMap}}
	if err := validateCmdArgs(CmdCapture, argMap); err != nil {
		return nil, err
	}

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

// GetCliString returns the raw CLI command string
func (c *CaptureCommand) GetCmdString() string {
	return c.cmd.args["cmd"]
}

// GetParsedCli returns the CLI command name to be captured and its arguments
func (c *CaptureCommand) GetParsedCmd() (string, []string) {
	return cliParse(c.cmd.args["cmd"])
}

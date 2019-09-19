// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

// CaptuerCommand represents CAPTURE in a script
type CaptureCommand struct {
	cmd
	cliCmdName string
	cliCmdArgs []string
}

// NewCaptureCommand returns *CaptureCommand with parsed arguments
func NewCaptureCommand(index int, args []string) (*CaptureCommand, error) {
	cmd := &CaptureCommand{cmd: cmd{index: index, name: CmdCapture, args: args}}

	if err := validateCmdArgs(CmdCapture, args); err != nil {
		return nil, err
	}

	cmdName, cmdArgs := cliParse(cmd.args[0])
	cmd.cliCmdName = cmdName
	cmd.cliCmdArgs = cmdArgs

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
func (c *CaptureCommand) Args() []string {
	return c.cmd.args
}

// GetCliString returns the raw CLI command string
func (c *CaptureCommand) GetCliString() string {
	return c.args[0]
}

// GetParsedCli returns the CLI command name to be captured and its arguments
func (c *CaptureCommand) GetParsedCli() (string, []string) {
	return c.cliCmdName, c.cliCmdArgs
}

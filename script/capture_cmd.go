// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

type CaptureCommand struct {
	cmd
	cliCmdName string
	cliCmdArgs []string
}

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

func (c *CaptureCommand) Index() int {
	return c.cmd.index
}

func (c *CaptureCommand) Name() string {
	return c.cmd.name
}

func (c *CaptureCommand) Args() []string {
	return c.cmd.args
}

func (c *CaptureCommand) GetCliString() string {
	return c.args[0]
}

func (c *CaptureCommand) GetParsedCli() (string, []string) {
	return c.cliCmdName, c.cliCmdArgs
}

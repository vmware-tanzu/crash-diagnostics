// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

// CopyCommand represents COPY command in a script
type CopyCommand struct {
	cmd
}

// NewCopyCommand returns *CopyCommand
func NewCopyCommand(index int, args []string) (*CopyCommand, error) {
	cmd := &CopyCommand{cmd: cmd{index: index, name: CmdCopy, args: args}}

	if err := validateCmdArgs(CmdCopy, args); err != nil {
		return nil, err
	}
	return cmd, nil
}

// Index is the position of the command in the script
func (c *CopyCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *CopyCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *CopyCommand) Args() []string {
	return c.cmd.args
}

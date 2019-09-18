// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

type CopyCommand struct {
	cmd
}

func NewCopyCommand(index int, args []string) (*CopyCommand, error) {
	cmd := &CopyCommand{cmd: cmd{index: index, name: CmdCopy, args: args}}

	if err := validateCmdArgs(CmdCopy, args); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *CopyCommand) Index() int {
	return c.cmd.index
}

func (c *CopyCommand) Name() string {
	return c.cmd.name
}

func (c *CopyCommand) Args() []string {
	return c.cmd.args
}

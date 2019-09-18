// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

type WorkdirCommand struct {
	cmd
	dir string
}

func NewWorkdirCommand(index int, args []string) (*WorkdirCommand, error) {
	cmd := &WorkdirCommand{cmd: cmd{index: index, name: CmdWorkDir, args: args}}

	if err := validateCmdArgs(cmd.name, args); err != nil {
		return nil, err
	}
	for _, arg := range args {
		cmd.dir = arg
		break // only get first arg.
	}
	return cmd, nil
}

func (c *WorkdirCommand) Index() int {
	return c.cmd.index
}

func (c *WorkdirCommand) Name() string {
	return c.cmd.name
}

func (c *WorkdirCommand) Args() []string {
	return c.cmd.args
}

func (c *WorkdirCommand) Dir() string {
	return c.dir
}

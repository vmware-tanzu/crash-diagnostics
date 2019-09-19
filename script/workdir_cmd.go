// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

// WorkdirCommand representes a WORKDIR directive in as script
type WorkdirCommand struct {
	cmd
	dir string
}

// NewWorkdirCommand parses args and returns a new *WorkdirCommand value
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

// Index is the position of the command in the script
func (c *WorkdirCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *WorkdirCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *WorkdirCommand) Args() []string {
	return c.cmd.args
}

// Dir returns the parsed path for directory
func (c *WorkdirCommand) Dir() string {
	return c.dir
}

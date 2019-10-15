// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

// CopyCommand represents a COPY directive:
//
// COPY [paths:"]path0 path1 pathN["]
type CopyCommand struct {
	cmd
}

// NewCopyCommand returns *CopyCommand
func NewCopyCommand(index int, rawArgs string) (*CopyCommand, error) {
	if err := validateRawArgs(CmdCopy, rawArgs); err != nil {
		return nil, err
	}

	// by default the args are assumed to be the path-list
	argMap := map[string]string{"paths": rawArgs}
	cmd := &CopyCommand{cmd: cmd{index: index, name: CmdCopy, args: argMap}}
	if err := validateCmdArgs(CmdCopy, argMap); err != nil {
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

// Paths returned a []string of paths to copy
func (c *CopyCommand) Paths() []string {
	return spaceSep.Split(c.cmd.args["paths"], -1)
}

// Args returns a slice of raw command arguments
func (c *CopyCommand) Args() map[string]string {
	return c.cmd.args
}

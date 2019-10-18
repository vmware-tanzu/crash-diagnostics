// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
)

// WorkdirCommand representes a WORKDIR which may have one
// of the following forms:
//
//    WORKDIR /path/to/workdir
//    WORKDIR path:/path/to/workdir
type WorkdirCommand struct {
	cmd
}

// NewWorkdirCommand parses args and returns a new *WorkdirCommand value
func NewWorkdirCommand(index int, rawArgs string) (*WorkdirCommand, error) {
	if err := validateRawArgs(CmdOutput, rawArgs); err != nil {
		return nil, err
	}

	var argMap map[string]string
	if !isNamedParam(rawArgs) {
		// setup default param (notice quoted value)
		rawArgs = makeNamedPram("path", rawArgs)
	}
	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("WORKDIR: %v", err)
	}

	cmd := &WorkdirCommand{cmd: cmd{index: index, name: CmdWorkDir, args: argMap}}
	if err := validateCmdArgs(cmd.name, argMap); err != nil {
		return nil, err
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
func (c *WorkdirCommand) Args() map[string]string {
	return c.cmd.args
}

// Path returns the parsed path for directory
func (c *WorkdirCommand) Path() string {
	return c.cmd.args["path"]
}

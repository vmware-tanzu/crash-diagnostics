// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
)

// OutputCommand representes a OUTPUT directive which can have
// one of the following forms:
//     OUTPUT /path/to/output
//     OUTPUT path:/path/to/output
type OutputCommand struct {
	cmd
}

// NewOutputCommand parses args and returns a new *OutputCommand value
func NewOutputCommand(index int, rawArgs string) (*OutputCommand, error) {
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
		return nil, fmt.Errorf("OUTPUT: %v", err)
	}

	cmd := &OutputCommand{cmd: cmd{index: index, name: CmdOutput, args: argMap}}
	if err := validateCmdArgs(cmd.name, argMap); err != nil {
		return nil, err
	}
	return cmd, nil
}

// Index is the position of the command in the script
func (c *OutputCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *OutputCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *OutputCommand) Args() map[string]string {
	return c.cmd.args
}

// Path returns the parsed path for directory
func (c *OutputCommand) Path() string {
	return ExpandEnv(c.cmd.args["path"])
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"strings"
)

// OutputCommand representes a OUTPUT directive in as script as
// OUTPUT path:<output-dir>
type OutputCommand struct {
	cmd
	path string
}

// NewOutputCommand parses args and returns a new *OutputCommand value
func NewOutputCommand(index int, args []string) (*OutputCommand, error) {
	cmd := &OutputCommand{cmd: cmd{index: index, name: CmdOutput, args: args}}

	if err := validateCmdArgs(cmd.name, args); err != nil {
		return nil, err
	}

	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s: bad argument: %s", cmd.name, arg)
		}
		switch {
		case strings.EqualFold(parts[0], "path"):
			cmd.path = parts[1]
		default:
		}
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
func (c *OutputCommand) Args() []string {
	return c.cmd.args
}

// Dir returns the parsed path for directory
func (c *OutputCommand) Path() string {
	return c.path
}

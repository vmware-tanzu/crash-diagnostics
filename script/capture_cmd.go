// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
)

// CaptureCommand represents CAPTURE directive which
// can have one of the following two forms as shown below:
//
//     CAPTURE <command-string>
//     CAPTURE cmd:"<command-string>" name:"cmd-name" desc:"cmd-desc"
//
// The former takes no named parameter. When the latter form is used,
// parameter cmd: is required.
type CaptureCommand struct {
	*RunCommand
}

// NewCaptureCommand returns *CaptureCommand with parsed arguments
func NewCaptureCommand(index int, rawArgs string) (*CaptureCommand, error) {
	runCmd, err := NewRunCommand(index, rawArgs)
	if err != nil {
		return nil, fmt.Errorf("CAPTURE: %v", err)
	}
	runCmd.name = CmdCapture

	return &CaptureCommand{runCmd}, nil
}

// GetEffectiveCmd returns the shell (if any) and command as
// a slice of strings
func (c *CaptureCommand) GetEffectiveCmd() ([]string, error) {
	args, err := c.RunCommand.GetEffectiveCmd()
	if err != nil {
		return nil, fmt.Errorf("CAPTURE: %s", err)
	}
	return args, nil
}

// GetParsedCmd returns the effective parsed command as commandName
// followed by a slice of command arguments
func (c *CaptureCommand) GetParsedCmd() (string, []string, error) {
	cmd, args, err := c.RunCommand.GetParsedCmd()
	if err != nil {
		return "", nil, fmt.Errorf("CAPTURE: %s", err)
	}
	return cmd, args, nil
}

// GetEffectiveCmdStr returns the effective command as a string
// which wraps the command around a shell quote if necessary
func (c *CaptureCommand) GetEffectiveCmdStr() (string, error) {
	cmdStr, err := c.RunCommand.GetEffectiveCmdStr()
	if err != nil {
		return "", fmt.Errorf("CAPTURE: %s", err)
	}
	return cmdStr, nil
}

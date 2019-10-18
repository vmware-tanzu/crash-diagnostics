// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
)

// KubeConfigCommand represents a KUBECONFIG directive which can have
// one of the following forms:
//     KUBECONFIG path/to/kubeconfig
//     KUBECONFIG path:"path/to/kubeconfig"
type KubeConfigCommand struct {
	cmd
}

// NewKubeConfigCommand creates a value of type KubeConfigCommand in a script
func NewKubeConfigCommand(index int, rawArgs string) (*KubeConfigCommand, error) {
	if err := validateRawArgs(CmdKubeConfig, rawArgs); err != nil {
		return nil, err
	}

	var argMap map[string]string
	if !isNamedParam(rawArgs) {
		// setup default param (notice quoted value)
		rawArgs = makeNamedPram("path", rawArgs)
	}
	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("KUBECONFIG: %v", err)
	}

	cmd := &KubeConfigCommand{cmd: cmd{index: index, name: CmdKubeConfig, args: argMap}}
	if err := validateCmdArgs(CmdKubeConfig, argMap); err != nil {
		return nil, err
	}
	cmd.cmd.args["path"] = searchForConfig(argMap["path"])
	return cmd, nil
}

// Index is the position of the command in the script
func (c *KubeConfigCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *KubeConfigCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *KubeConfigCommand) Args() map[string]string {
	return c.cmd.args
}

// Config returns the path to the config file
func (c *KubeConfigCommand) Path() string {
	return c.cmd.args["path"]
}

// searchForConfig searches in several places for
// the kubernets config:
// 1. from passed args
// 2. from ENV variable
// 3. from local homedir
func searchForConfig(defaultPath string) string {
	if len(defaultPath) > 0 {
		return defaultPath
	}
	return Defaults.KubeConfigValue
}

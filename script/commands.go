// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"path/filepath"
)

var (
	CmdAs          = "AS"
	CmdCapture     = "CAPTURE"
	CmdCopy        = "COPY"
	CmdEnv         = "ENV"
	CmdFrom        = "FROM"
	CmdKubeConfig  = "KUBECONFIG"
	CmdSSHConfig   = "SSHCONFIG"
	CmdFromDefault = "local"
	CmdWorkDir     = "WORKDIR"

	Defaults = struct {
		FromValue       string
		WorkdirValue    string
		KubeConfigValue string
	}{
		FromValue:    "local",
		WorkdirValue: "/tmp/flareout",
		KubeConfigValue: func() string {
			kubecfg := os.Getenv("KUBECONFIG")
			if kubecfg == "" {
				kubecfg = filepath.Join(os.Getenv("HOME"), ".kube", "config")
			}
			return kubecfg
		}(),
	}
)

type CommandMeta struct {
	Name      string
	MinArgs   int
	MaxArgs   int
	Supported bool
}

var (
	Cmds = map[string]CommandMeta{
		CmdAs:         CommandMeta{Name: CmdAs, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdCapture:    CommandMeta{Name: CmdCapture, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdCopy:       CommandMeta{Name: CmdCopy, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdEnv:        CommandMeta{Name: CmdEnv, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdFrom:       CommandMeta{Name: CmdFrom, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdKubeConfig: CommandMeta{Name: CmdKubeConfig, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdSSHConfig:  CommandMeta{Name: CmdSSHConfig, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdWorkDir:    CommandMeta{Name: CmdWorkDir, MinArgs: 1, MaxArgs: 1, Supported: true},
	}
)

// Command is an abtract representatio of command in a script
type Command interface {
	// Index is the position of the command in the script
	Index() int
	// Name represents the name of the command
	Name() string
	// Args returns a slice of raw command arguments
	Args() []string
}

// Script is a collection of commands
type Script struct {
	Preambles map[string][]Command // directive commands in the script
	Actions   []Command            // action commands
}

// cmd is the base representation of command
type cmd struct {
	index int
	name  string
	args  []string
}

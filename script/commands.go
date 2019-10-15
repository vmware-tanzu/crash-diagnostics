// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"path/filepath"
)

var (
	CmdAs         = "AS"
	CmdAuthConfig = "AUTHCONFIG"
	CmdCapture    = "CAPTURE"
	CmdCopy       = "COPY"
	CmdEnv        = "ENV"
	CmdFrom       = "FROM"
	CmdKubeConfig = "KUBECONFIG"
	CmdOutput     = "OUTPUT"
	CmdWorkDir    = "WORKDIR"

	Defaults = struct {
		FromValue       string
		WorkdirValue    string
		KubeConfigValue string
		AuthPKValue     string
		OutputValue     string
	}{
		FromValue:    "local",
		WorkdirValue: "/tmp/crashdir",
		KubeConfigValue: func() string {
			kubecfg := os.Getenv("KUBECONFIG")
			if kubecfg == "" {
				kubecfg = filepath.Join(os.Getenv("HOME"), ".kube", "config")
			}
			return kubecfg
		}(),
		AuthPKValue: func() string {
			return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
		}(),
		OutputValue: "out.tar.gz",
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
		CmdAs:         CommandMeta{Name: CmdAs, MinArgs: 1, MaxArgs: 2, Supported: true},
		CmdAuthConfig: CommandMeta{Name: CmdAuthConfig, MinArgs: 1, MaxArgs: 3, Supported: true},
		CmdCapture:    CommandMeta{Name: CmdCapture, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdCopy:       CommandMeta{Name: CmdCopy, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdEnv:        CommandMeta{Name: CmdEnv, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdFrom:       CommandMeta{Name: CmdFrom, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdKubeConfig: CommandMeta{Name: CmdKubeConfig, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdOutput:     CommandMeta{Name: CmdOutput, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdWorkDir:    CommandMeta{Name: CmdWorkDir, MinArgs: 1, MaxArgs: 1, Supported: true},
	}
)

type ArgMap = map[string]string

// Command is an abtract representatio of command in a script
type Command interface {
	// Index is the position of the command in the script
	Index() int
	// Name represents the name of the command
	Name() string
	// Args returns a map of parsed arguments
	Args() ArgMap
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
	args  map[string]string
}

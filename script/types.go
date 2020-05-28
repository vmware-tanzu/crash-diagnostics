// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"path/filepath"
)

var (
	CmdAs         = "AS"
	CmdAsConfig   = "ASCONFIG"
	CmdAuthConfig = "AUTHCONFIG"
	CmdCapture    = "CAPTURE"
	CmdCopy       = "COPY"
	CmdEnv        = "ENV"
	CmdFrom       = "FROM"
	CmdKubeConfig = "KUBECONFIG"
	CmdKubeGet    = "KUBEGET"
	CmdOutput     = "OUTPUT"
	CmdRun        = "RUN"
	CmdWorkDir    = "WORKDIR"

	Defaults = struct {
		FromValue         string
		WorkdirValue      string
		KubeConfigValue   string
		AuthPKValue       string
		OutputValue       string
		HostAddr          string
		ServicePort       string
		ConnectionRetries string
		ConnectionTimeout string
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
		OutputValue:       "out.tar.gz",
		HostAddr:          "127.0.0.1",
		ServicePort:       "22",
		ConnectionRetries: "30",
		ConnectionTimeout: "2m",
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
		CmdCapture:    CommandMeta{Name: CmdCapture, MinArgs: 1, MaxArgs: 3, Supported: true},
		CmdCopy:       CommandMeta{Name: CmdCopy, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdEnv:        CommandMeta{Name: CmdEnv, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdFrom:       CommandMeta{Name: CmdFrom, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdKubeConfig: CommandMeta{Name: CmdKubeConfig, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdKubeGet:    CommandMeta{Name: CmdKubeGet, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdOutput:     CommandMeta{Name: CmdOutput, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdRun:        CommandMeta{Name: CmdRun, MinArgs: 1, MaxArgs: 3, Supported: true},
		CmdWorkDir:    CommandMeta{Name: CmdWorkDir, MinArgs: 1, MaxArgs: 1, Supported: true},
	}
)

type ArgMap = map[string]string

// Directive base interface that represents a directive
// Implementation should provide capture ample parameters
// so the directive is properly handled at runtime.
type Directive interface {
	// Index position of the command in the script
	Index() int
	// Name the raw name of the command
	Name() string

	Raw() string
}

// ConfigDirective marker interface that represents a configuration
type ConfigDirective interface {
	Directive
}

// ExecDirective marker interface for an executable directive
type ExecDirective interface {
	Directive
}

// cmd is the base representation of command
type cmd struct {
	index int
	name  string
	args  map[string]string
}

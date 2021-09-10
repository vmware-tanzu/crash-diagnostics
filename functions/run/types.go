// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package run implements starlark function 'run()`
// used to execute processes on remote compute resources.
package run

import (
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
)

type Args struct {
	Cmd       string              `name:"cmd"`
	SSHConfig sshconf.Config      `name:"ssh_config" optional:"true"`
	Resources providers.Resources `name:"resources" optional:"true"`
}

type RemoteProc struct {
	Error  string `name:"error"`
	Host   string `name:"host"`
	Output string `name:"output"`
}

type Result struct {
	Error string       `name:"error"`
	Procs []RemoteProc `name:"procs"`
}

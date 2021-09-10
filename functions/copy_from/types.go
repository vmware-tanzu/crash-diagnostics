// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package copy_from represents the `copy_from` Starlark function.
package copy_from

import (
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
)

type Args struct {
	Path      string              `name:"path"`
	SSHConfig sshconf.Config      `name:"ssh_config" optional:"true"`
	Resources providers.Resources `name:"resources" optional:"true"`
	Workdir   string              `name:"workdir" optional:"true"`
}

type RemoteCopy struct {
	Error string `name:"error"`
	Host  string `name:"host"`
	Path  string `name:"path"`
}

type Result struct {
	Error  string       `name:"error"`
	Copies []RemoteCopy `name:"copies"`
}

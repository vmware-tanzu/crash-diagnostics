// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package copy_to represents the `copy_to` Starlark function.
package copy_to

import (
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
)

type Args struct {
	SourcePath string              `name:"source_path"`
	TargetPath string              `name:"target_path" optional:"true"`
	SSHConfig  sshconf.Config      `name:"ssh_config" optional:"true"`
	Resources  providers.Resources `name:"resources" optional:"true"`
}

type CopyResult struct {
	Error      string `name:"error"`
	Host       string `name:"host"`
	SourcePath string `name:"source_path"'`
	TargetPath string `name:"target_path"`
}

type Result struct {
	Error  string       `name:"error"`
	Copies []CopyResult `name:"copies"`
}

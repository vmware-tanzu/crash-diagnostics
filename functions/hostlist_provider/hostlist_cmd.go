// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package hostlist_provider

import (
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"go.starlark.net/starlark"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

func (c *cmd) Run(t *starlark.Thread, args Args) functions.ProviderResources {
	if len(args.Hosts) == 0 {
		return functions.ProviderResources{Error: "host list is required"}
	}

	return functions.ProviderResources{
		Hosts: args.Hosts,
	}
}
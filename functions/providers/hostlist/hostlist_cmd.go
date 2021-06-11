// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package hostlist

import (
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"go.starlark.net/starlark"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

func (c *cmd) Run(t *starlark.Thread, args Args) providers.Resources {
	if len(args.Hosts) == 0 {
		return providers.Resources{Error: "host list is required"}
	}

	return providers.Resources{
		Provider: string(Name),
		Hosts:    args.Hosts,
	}
}

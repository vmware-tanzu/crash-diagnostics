// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package copy_from

import (
	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
	"go.starlark.net/starlark"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

// TODO complete
func (c *cmd) Run(t *starlark.Thread, args Args) Result {
	if args.Path == "" {
		return Result{Error: "no source path provided"}
	}

	if args.Workdir == "" {
		if conf, ok := scriptconf.ConfigFromThread(t); ok {
			args.Workdir = conf.Workdir
		} else {
			args.Workdir = scriptconf.DefaultWorkdir()
		}
	}
	return Result{}
}

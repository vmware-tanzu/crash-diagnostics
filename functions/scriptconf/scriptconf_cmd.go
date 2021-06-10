// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"fmt"
	"os"
	"os/user"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"github.com/vmware-tanzu/crash-diagnostics/util"
)

type confCmd struct{}

func newCmd() *confCmd {
	return new(confCmd)
}

// Run applies processes the params and generates a configuration value for the script
func (c *confCmd) Run(t *starlark.Thread, args Args) Result {
	if err := validateArgs(&args); err != nil {
		return Result{Error: fmt.Sprintf("failed to validate configuration: %s", err)}
	}

	// create workdir if needed
	if err := functions.MakeDir(args.Workdir, 0744); err != nil && !os.IsExist(err) {
		return Result{Error: fmt.Sprintf("failed to create workdir: %s", err)}
	}

	// start local ssh-agent
	if args.UseSSHAgent {
		agent, err := ssh.StartAgent()
		if err != nil {
			return Result{Error: errors.Wrap(err, "failed to start ssh agent").Error()}
		}
		t.SetLocal(sshconf.SSHAgentIdentifier, agent)
	}

	return Result{
		Workdir:      args.Workdir,
		Gid:          args.Gid,
		Uid:          args.Uid,
		DefaultShell: args.DefaultShell,
		Requires:     args.Requires,
		UseSSHAgent:  args.UseSSHAgent,
	}
}

func validateArgs(params *Args) error {
	if params.Workdir == "" {
		params.Workdir = DefaultWorkdir()
	}
	wd, err := util.ExpandPath(params.Workdir)
	if err != nil {
		return err
	}
	params.Workdir = wd

	if params.Gid == "" {
		params.Gid = getGid()
	}

	if params.Uid == "" {
		params.Uid = getUid()
	}

	return nil
}

func getGid() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Gid
}

func getUid() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Uid
}

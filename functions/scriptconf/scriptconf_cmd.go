// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"fmt"
	"os/user"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"github.com/vmware-tanzu/crash-diagnostics/util"
)

const (
	defaultWorkdir = "/tmp/crashd"
)

type confCmd struct{}

func newCmd() *confCmd {
	return new(confCmd)
}

// Run applies processes the params and generates a configuration value for the script
func (c *confCmd) Run(t *starlark.Thread, p interface{}) (functions.CommandResult, error) {
	params, ok := p.(Params)
	if !ok {
		return nil, fmt.Errorf("unexpected params type: %T", p)
	}

	if err := validateParams(&params); err != nil {
		return nil, fmt.Errorf("failed to build configuration: %w", err)
	}

	if params.Workdir != "" {
		if err := functions.MakeDir(params.Workdir, 0744); err != nil {
			return nil, fmt.Errorf("failed to create workdir: %w", err)
		}
	}

	// start local ssh-agent
	if params.UseSSHAgent {
		agent, err := ssh.StartAgent()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ssh agent")
		}
		t.SetLocal("ssh_agent", agent)
	}

	return functions.NewResult(params), nil
}

func validateParams(params *Params) error {
	if params.Workdir == "" {
		params.Workdir = defaultWorkdir
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

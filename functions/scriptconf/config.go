// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"fmt"
	"os/user"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"github.com/vmware-tanzu/crash-diagnostics/util"
)

const (
	defaultWorkdir = "/tmp/crashd"
)

type Params struct {
	Workdir      string
	Gid          string
	Uid          string
	DefaultShell string
	Requires     []string
	UseSSHAgent  bool
}

// Configuration is a mirror of param (in this case)
type Configuration = Params

// Build creates the configuration value for the script
func Build(t *starlark.Thread, params Params) (Configuration, error) {

	if err := validateParams(&params); err != nil {
		return Params{}, fmt.Errorf("failed to build configuration: %w", err)
	}

	// start local ssh-agent
	if params.UseSSHAgent {
		agent, err := ssh.StartAgent()
		if err != nil {
			return Params{}, errors.Wrap(err, "failed to start ssh agent")
		}
		t.SetLocal("ssh_agent", agent)
	}

	return params, nil
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

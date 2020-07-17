// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// SetAsDefaultFunc is the built-in fn that saves the arguments to the local Starlark thread.
// Starlark format: set_as_default([ssh_config = ssh_config()][, kube_config = kube_config()][, resources = resources()])
func SetAsDefaultFunc(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var kubeConfig, sshConfig *starlarkstruct.Struct
	var resources *starlark.List

	if err := starlark.UnpackArgs(
		identifiers.setAsDefault, args, kwargs,
		"kube_config?", &kubeConfig,
		"ssh_config?", &sshConfig,
		"resources?", &resources,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.setAsDefault, err)
	}

	if sshConfig == nil && kubeConfig == nil && resources == nil {
		return starlark.None, errors.New("atleast one of kube_config, ssh_config or resources is required")
	}

	if kubeConfig != nil {
		thread.SetLocal(identifiers.kubeCfg, kubeConfig)
	}
	if sshConfig != nil {
		thread.SetLocal(identifiers.sshCfg, sshConfig)
	}
	if resources != nil {
		thread.SetLocal(identifiers.resources, resources)
	}

	return starlark.None, nil
}

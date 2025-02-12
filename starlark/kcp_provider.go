// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// KcpProviderFn is a built-in starlark function that create a kubeconfig with all contexts for all KCP logical clusters
// Starlark format: kcp_provider(ucp_admin_secret_name=<ucp_admin_secret_name> ucp_admin_secret_namespace=<ucp_admin_secret_namespace>)
func KcpProviderFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	var (
		kcp_admin_secret_name, kcp_admin_secret_namespace string
	)

	err := starlark.UnpackArgs("kcp_provider", args, kwargs,
		"kcp_admin_secret_namespace", &kcp_admin_secret_namespace,
		"kcp_admin_secret_name", &kcp_admin_secret_name)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to unpack input arguments")
	}

	ctx, ok := thread.Local(identifiers.scriptCtx).(context.Context)
	if !ok || ctx == nil {
		return starlark.None, fmt.Errorf("script context not found")
	}

	var kcpKubeConfigPath = "/Users/tatanas/dev.kubeconfig"
	//TODO Generate a KCP admin kubeconfig

	// dictionary for capa provider struct
	kcpProviderDict := starlark.StringDict{
		"kind":        starlark.String(identifiers.kcpProvider),
		"kube_config": starlark.String(kcpKubeConfigPath),
	}

	var contexts []starlark.Value
	contexts = append(contexts, starlark.String("tanzu-cli-Falcons_GCP_New-staging-5a2f0150:project-ashindov"))
	contexts = append(contexts, starlark.String("tanzu-cli-Falcons_GCP_New-staging-5a2f0150:project-ashindov:services-space"))

	kcpProviderDict["contexts"] = starlark.NewList(contexts)

	return starlarkstruct.FromStringDict(starlark.String(identifiers.kcpProvider), kcpProviderDict), nil
}

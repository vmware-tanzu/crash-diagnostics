// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"context"
	"fmt"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

// UcpProviderFn is a built-in starlark function that collects Kubconfigs for all available UCP workspaces
// Starlark format: ucp_provider(kube_config=kube_config())
func UcpProviderFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	var (
		mgmtKubeConfig *starlarkstruct.Struct
	)

	err := starlark.UnpackArgs("ucp_provider", args, kwargs,
		"mgmt_kube_config", &mgmtKubeConfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to unpack input arguments")
	}

	ctx, ok := thread.Local(identifiers.scriptCtx).(context.Context)
	if !ok || ctx == nil {
		return starlark.None, fmt.Errorf("script context not found")
	}

	if mgmtKubeConfig == nil {
		mgmtKubeConfig = thread.Local(identifiers.kubeCfg).(*starlarkstruct.Struct)
	}

	ucpProviderDict := starlark.StringDict{
		"kind":             starlark.String(identifiers.ucpProvider),
		"transport":        starlark.String("ssh"),
		identifiers.sshCfg: starlark.String("sshCfg"),
	}
	mgmtKubeConfigPath, err := getKubeConfigPathFromStruct(mgmtKubeConfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to extract management kubeconfig")
	}

	config, err := clientcmd.LoadFromFile(mgmtKubeConfigPath)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Step 3: Create a KCP client using the kubeconfig
	kcpConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: mgmtKubeConfigPath},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		log.Fatalf("Failed to create KCP client config: %v", err)
	}

	kcpClient, err := kcpclient.NewForConfig(kcpConfig)
	if err != nil {
		log.Fatalf("Failed to create KCP client: %v", err)
	}

	workspaces, err := kcpClient.TenancyV1alpha1().Workspaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list KCP workspaces: %v", err)
	}

	var kubeconfigs []starlark.Value
	for _, ws := range workspaces.Items {
		switchWorkspace := NewUseWorkspaceCommand(ws.Name, config, mgmtKubeConfigPath)
		kubeConfigPath, err := switchWorkspace.Run(context.Background())
		if err != nil {
			return nil, err
		}
		kubeconfigs = append(kubeconfigs, starlark.String(kubeConfigPath))
	}
	ucpProviderDict["hosts"] = starlark.NewList(kubeconfigs)

	return starlarkstruct.FromStringDict(starlark.String(identifiers.ucpProvider), ucpProviderDict), nil
}

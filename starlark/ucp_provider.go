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
// Starlark format: ucp_provider(kube_config=kube_config(),workspace=<name>)
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
	//mgmtKubeConfigPath, err := getKubeConfigPathFromStruct(mgmtKubeConfig)
	//if err != nil {
	//	return starlark.None, errors.Wrap(err, "failed to extract management kubeconfig")
	//}

	// if workload cluster is not supplied, then the resources for the management cluster
	// should be enumerated
	//workspaceName := workspace
	//if workspaceName == "" {
	//	//config, err := k8s.LoadKubeCfg(mgmtKubeConfigPath)
	//	if err != nil {
	//		return starlark.None, errors.Wrap(err, "failed to load kube config")
	//	}
	//	workspaceName, err = config.GetClusterName()
	//	if err != nil {
	//		return starlark.None, errors.Wrap(err, "cannot find cluster with name "+wo)
	//	}
	//}
	//
	//providerConfigPath, err := provider.KubeConfig(mgmtKubeConfigPath, clusterName, namespace)
	//if err != nil {
	//	return starlark.None, err
	//}
	//
	//nodeAddresses, err := k8s.GetNodeAddresses(ctx, providerConfigPath, toSlice(names), toSlice(labels))
	//if err != nil {
	//	return starlark.None, errors.Wrap(err, "could not fetch host addresses")
	//}

	//dictionary for ucp provider struct
	ucpProviderDict := starlark.StringDict{
		"kind":             starlark.String(identifiers.ucpProvider),
		"transport":        starlark.String("ssh"),
		identifiers.sshCfg: starlark.String("sshCfg"),
		//"kube_config": starlark.String(providerConfigPath),
	}
	mgmtKubeConfigPath, err := getKubeConfigPathFromStruct(mgmtKubeConfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to extract management kubeconfig")
	}
	//
	//loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	//loadingRules.ExplicitPath = mgmtKubeConfigPath
	//
	//startingConfig, err := loadingRules.GetStartingConfig()

	config, err := clientcmd.LoadFromFile(mgmtKubeConfigPath)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	//switchWorkspace := NewUseWorkspaceCommand(workspace, config, mgmtKubeConfigPath)
	//err = switchWorkspace.Run(context.Background())
	//if err != nil {
	//	return nil, err
	//}

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

	//// Find the cluster referenced by the current context
	//kubeContext := config.Contexts[config.CurrentContext]
	//clusterName := kubeContext.Cluster
	//cluster := config.Clusters["workspace.kcp.io/current"]
	//if cluster == nil {
	//	fmt.Printf("No cluster named %q found in kubeconfig.\n", clusterName)
	//	return nil, errors.Wrap(err, "No cluster named %"+clusterName)
	//}

	//currentClusterUrl, err := url.Parse(cluster.Server)
	//if err != nil {
	//	panic(err)
	//}

	//cluster.Server += ":space-1"

	//https://localhost:6453/clusters/root:testorg:project-hr

	//err = clientcmd.WriteToFile(*config, mgmtKubeConfigPath)
	//if err != nil {
	//	log.Fatalf("Failed to save updated kubeconfig: %v", err)
	//}

	wcs := make([]starlark.Value, 0)
	var kubeconfigs []starlark.Value
	for _, ws := range workspaces.Items {
		switchWorkspace := NewUseWorkspaceCommand(ws.Name, config, mgmtKubeConfigPath)
		kubeConfigPath, err := switchWorkspace.Run(context.Background())
		if err != nil {
			return nil, err
		}
		wcs = append(wcs, starlark.String(kubeConfigPath))
		kubeconfigs = append(kubeconfigs, starlark.String(kubeConfigPath))
	}

	ucpProviderDict["kubeconfigs"] = starlark.NewList(kubeconfigs)
	//kubconfigs := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	//// add node info to dictionary
	//var nodeIps []starlark.Value
	//for _, kc := range kubconfigs {
	//	nodeIps = append(nodeIps, starlark.String(kc))
	//}
	ucpProviderDict["hosts"] = starlark.NewList(wcs)

	//sshConfigDict := starlark.StringDict{}
	//sshConfig.ToStringDict(sshConfigDict)
	//
	//// modify ssh config jump credentials, if not specified
	//if _, err := sshConfig.Attr("jump_host"); err != nil {
	//	sshConfigDict["jump_host"] = starlark.String(bastionIpAddr)
	//}
	//if _, err := sshConfig.Attr("jump_user"); err != nil {
	//	sshConfigDict["jump_user"] = starlark.String("ubuntu")
	//}
	//capaProviderDict[identifiers.sshCfg] = starlarkstruct.FromStringDict(starlark.String(identifiers.sshCfg), sshConfigDict)

	return starlarkstruct.FromStringDict(starlark.String(identifiers.ucpProvider), ucpProviderDict), nil
	//return nil, errors.Wrap(err, "Not yet implemented")
}

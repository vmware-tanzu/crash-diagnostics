// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/provider"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// CapvProviderFn is a built-in starlark function that collects compute resources from a k8s cluster
// Starlark format: capv_provider(kube_config=kube_config(), ssh_config=ssh_config()[workload_cluster=<name>, nodes=["foo", "bar], labels=["bar", "baz"]])
func CapvProviderFn(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	var (
		workloadCluster       string
		names, labels         *starlark.List
		sshConfig, kubeConfig *starlarkstruct.Struct
	)

	err := starlark.UnpackArgs("capv_provider", args, kwargs,
		"ssh_config", &sshConfig,
		"kube_config", &kubeConfig,
		"workload_cluster?", &workloadCluster,
		"labels?", &labels,
		"nodes?", &names)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to unpack input arguments")
	}

	if sshConfig == nil || kubeConfig == nil {
		return starlark.None, errors.New("capv_provider requires the name of the management cluster, the ssh configuration and the management cluster kubeconfig")
	}

	mgmtKubeConfigPath, err := getKubeConfigFromStruct(kubeConfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to extract management kubeconfig")
	}

	providerConfigPath, err := provider.KubeConfig(mgmtKubeConfigPath, workloadCluster)
	if err != nil {
		return starlark.None, err
	}

	nodeAddresses, err := k8s.GetNodeAddresses(providerConfigPath, toSlice(names), toSlice(labels))
	if err != nil {
		return starlark.None, errors.Wrap(err, "could not fetch host addresses")
	}

	// dictionary for capv provider struct
	capvProviderDict := starlark.StringDict{
		"kind":       starlark.String(identifiers.capvProvider),
		"transport":  starlark.String("ssh"),
		"kubeconfig": starlark.String(providerConfigPath),
	}

	// add node info to dictionary
	var nodeIps []starlark.Value
	for _, node := range nodeAddresses {
		nodeIps = append(nodeIps, starlark.String(node))
	}
	capvProviderDict["hosts"] = starlark.NewList(nodeIps)

	// add ssh info to dictionary
	if _, ok := capvProviderDict[identifiers.sshCfg]; !ok {
		capvProviderDict[identifiers.sshCfg] = sshConfig
	}

	return starlarkstruct.FromStringDict(starlark.String(identifiers.capvProvider), capvProviderDict), nil
}

// TODO: Needs to be moved to a single package
func toSlice(list *starlark.List) []string {
	var elems []string
	if list != nil {
		for i := 0; i < list.Len(); i++ {
			if val, ok := list.Index(i).(starlark.String); ok {
				elems = append(elems, string(val))
			}
		}
	}
	return elems
}

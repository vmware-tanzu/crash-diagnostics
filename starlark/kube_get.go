// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// KubeGetFn is a starlark built-in for the fetching kubernetes objects
func KubeGetFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var objects *starlark.List
	var groups, kinds, namespaces, versions, names, labels, containers *starlark.List
	var kubeConfig *starlarkstruct.Struct

	if err := starlark.UnpackArgs(
		identifiers.crashdCfg, args, kwargs,
		"groups?", &groups,
		"kinds?", &kinds,
		"namespaces?", &namespaces,
		"versions?", &versions,
		"names?", &names,
		"labels?", &labels,
		"containers?", &containers,
		"kube_config?", &kubeConfig,
	); err != nil {
		return starlark.None, errors.Wrap(err, "failed to read args")
	}

	if kubeConfig == nil {
		kubeConfig = thread.Local(identifiers.kubeCfg).(*starlarkstruct.Struct)
	}
	path, err := getKubeConfigFromStruct(kubeConfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to kubeconfig")
	}

	client, err := k8s.New(path)
	if err != nil {
		return starlark.None, errors.Wrap(err, "could not initialize search client")
	}

	searchParams := k8s.SearchParams{
		Groups:     toSlice(groups),
		Kinds:      toSlice(kinds),
		Namespaces: toSlice(namespaces),
		Versions:   toSlice(versions),
		Names:      toSlice(names),
		Labels:     toSlice(labels),
		Containers: toSlice(containers),
	}
	searchResults, err := client.Search(searchParams)
	if err == nil {
		objects = starlark.NewList([]starlark.Value{})
		for _, searchResult := range searchResults {
			srValue := searchResult.ToStarlarkValue()
			err = objects.Append(srValue)
			if err != nil {
				err = errors.Wrap(err, "could not collect kube_get() results")
				break
			}
		}
	}

	return starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		starlark.StringDict{
			"objs": objects,
			"error": func() starlark.String {
				if err != nil {
					return starlark.String(err.Error())
				}
				return ""
			}(),
		}), nil
}

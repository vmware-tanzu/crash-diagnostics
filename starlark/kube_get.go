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
func KubeGetFn(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var objects *starlark.List

	structVal, err := kwargsToStruct(kwargs)
	if err != nil {
		return starlark.None, err
	}

	kubeconfig, err := getKubeConfigPath(thread, structVal)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to kubeconfig")
	}
	client, err := k8s.New(kubeconfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "could not initialize search client")
	}

	searchParams := k8s.NewSearchParams(structVal)
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

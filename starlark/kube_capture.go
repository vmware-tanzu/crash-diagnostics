// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// KubeCaptureFn is the Starlark built-in for the fetching kubernetes objects
// and returns the result as a Starlark value containing the file path and error message, if any
// Starlark format: kube_capture(what="logs" [, groups="core", namespaces=["default"], kube_config=kube_config()])
func KubeCaptureFn(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var argDict starlark.StringDict

	if kwargs != nil {
		dict, err := kwargsToStringDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		argDict = dict
	}
	structVal := starlarkstruct.FromStringDict(starlarkstruct.Default, argDict)

	kubeconfig, err := getKubeConfigPath(thread, structVal)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to kubeconfig")
	}
	client, err := k8s.New(kubeconfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "could not initialize search client")
	}

	data := thread.Local(identifiers.crashdCfg)
	cfg, _ := data.(*starlarkstruct.Struct)
	workDirVal, _ := cfg.Attr("workdir")
	resultDir, err := write(trimQuotes(workDirVal.String()), client, structVal)

	return starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		starlark.StringDict{
			"file": starlark.String(resultDir),
			"error": func() starlark.String {
				if err != nil {
					return starlark.String(err.Error())
				}
				return ""
			}(),
		}), nil
}

func write(workdir string, client *k8s.Client, structVal *starlarkstruct.Struct) (string, error) {
	var searchResults []k8s.SearchResult
	whatVal, err := structVal.Attr("what")
	// TODO: check if we need default value
	if err != nil {
		return "", errors.Wrap(err, "what input parameter not specified")
	}
	whatStrVal, _ := whatVal.(starlark.String)
	what := whatStrVal.GoString()

	searchParams := k8s.NewSearchParams(structVal)

	logrus.Debugf("kube_capture(what=%s)", what)
	switch what {
	case "logs":
		searchParams.Groups = []string{"core"}
		searchParams.Kinds = []string{"pods"}
		searchParams.Versions = []string{}
	case "objects", "all", "*":
	default:
		return "", errors.Errorf("don't know how to get: %s", what)
	}

	searchResults, err = client.Search(searchParams)
	if err != nil {
		return "", err
	}

	resultWriter, err := k8s.NewResultWriter(workdir, what, client.CoreRest)
	if err != nil {
		return "", errors.Wrap(err, "failed to initialize writer")
	}
	err = resultWriter.Write(searchResults)
	if err != nil {
		return "", errors.Wrap(err, "failed to write search results")
	}
	return resultWriter.GetResultDir(), nil
}

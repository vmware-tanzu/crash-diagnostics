// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeCaptureFn is the Starlark built-in for the fetching kubernetes objects
// and returns the result as a Starlark value containing the file path and error message, if any
// Starlark format: kube_capture(what="logs" [, groups="core", namespaces=["default"], kube_config=kube_config(), tunnel_config=tunnel_config])
func KubeCaptureFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	var groups, categories, kinds, namespaces, versions, names, labels, containers *starlark.List
	var kubeConfig *starlarkstruct.Struct
	var tunnelConfig *starlarkstruct.Struct
	var what string
	var outputFormat string
	var outputMode string
	logrus.Info(kwargs)

	if err := starlark.UnpackArgs(
		identifiers.kubeCapture, args, kwargs,
		"what", &what,
		"output_format?", &outputFormat,
		"output_mode?", &outputMode,
		"groups?", &groups,
		"categories?", &categories,
		"kinds?", &kinds,
		"namespaces?", &namespaces,
		"versions?", &versions,
		"names?", &names,
		"labels?", &labels,
		"containers?", &containers,
		"kube_config?", &kubeConfig,
		"tunnel_config?", &tunnelConfig,
	); err != nil {
		return starlark.None, fmt.Errorf("failed to read args: %w", err)
	}

	writeLogs := what == "logs" || what == "all"
	if writeLogs && tunnelConfig != nil {
		return starlark.None, fmt.Errorf("tunnel_config unsupported for 'logs' and 'all' operations")
	}

	ctx, ok := thread.Local(identifiers.scriptCtx).(context.Context)
	if !ok || ctx == nil {
		return starlark.None, errors.New("script context not found")
	}

	if kubeConfig == nil {
		kubeConfig = thread.Local(identifiers.kubeCfg).(*starlarkstruct.Struct)
	}
	path, err := getKubeConfigPathFromStruct(kubeConfig)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to kubeconfig: %w", err)
	}
	var targetClient *k8s.Client
	_, err = kubeConfig.Attr("extra_kubeconfig")
	usePortforward := err == nil
	clusterCtxName := getKubeConfigContextNameFromStruct(kubeConfig)

	if usePortforward {
		if targetClient, err = newTargetKubeconfig(kubeConfig, clusterCtxName); err != nil {
			return starlark.None, err
		}
	}

	path, err = getKubeConfigPathFromStruct(kubeConfig)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to kubeconfig: %w", err)
	}

	portforwardStr, err := tunnelConfig.Attr("namespace")
	if err != nil {
		return nil, fmt.Errorf("could not get set the roundtripper: %w", err)
	}

	portForwardNS := portforwardStr.(starlark.String)

	podNameStr, err := tunnelConfig.Attr("pod_name")
	if err != nil {
		return nil, fmt.Errorf("could not get set the roundtripper: %w", err)
	}

	portForwardPodName := podNameStr.(starlark.String)

	localPortInt, err := tunnelConfig.Attr("local_port")
	if err != nil {
		return nil, fmt.Errorf("could not get set the roundtripper: %w", err)
	}

	localPort := localPortInt.(starlark.Int)

	targetPortInt, err := tunnelConfig.Attr("target_port")
	if err != nil {
		return nil, fmt.Errorf("could not get set the roundtripper: %w", err)
	}
	targetPort := targetPortInt.(starlark.Int)

	fw, err := k8s.NewPortForwarder(path, string(portForwardNS), string(portForwardPodName), int(localPort.BigInt().Int64()), int(targetPort.BigInt().Int64()))
	if err != nil {
		return starlark.None, err
	}

	defer fw.Close()
	go func() error {
		if err := fw.ForwardPorts(); err != nil {
			return err
		}
		return nil
	}()

	data := thread.Local(identifiers.crashdCfg)
	cfg, _ := data.(*starlarkstruct.Struct)
	workDirVal, _ := cfg.Attr("workdir")
	resultDir, err := write(ctx, trimQuotes(workDirVal.String()), what, strings.ToLower(outputFormat), strings.ToLower(outputMode), targetClient, k8s.SearchParams{
		Groups:     toSlice(groups),
		Categories: toSlice(categories),
		Kinds:      toSlice(kinds),
		Namespaces: toSlice(namespaces),
		Versions:   toSlice(versions),
		Names:      toSlice(names),
		Labels:     toSlice(labels),
		Containers: toSlice(containers),
	})

	return starlarkstruct.FromStringDict(
		starlark.String(identifiers.kubeCapture),
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

func write(ctx context.Context, workdir, what, outputFormat, outputMode string, client *k8s.Client, params k8s.SearchParams) (string, error) {

	logrus.Debugf("kube_capture(what=%s)", what)
	switch what {
	case "logs":
		params.Groups = []string{"core"}
		params.Kinds = []string{"pods"}
		params.Versions = []string{}
	case "objects", "all", "*":
	default:
		return "", fmt.Errorf("don't know how to get: %s", what)
	}

	searchResults, err := client.Search(ctx, params)
	if err != nil {
		return "", err
	}

	resultWriter, err := k8s.NewResultWriter(workdir, what, outputFormat, outputMode, client.CoreRest)
	if err != nil {
		return "", fmt.Errorf("failed to initialize writer: %w", err)
	}
	err = resultWriter.Write(ctx, searchResults)
	if err != nil {
		return "", fmt.Errorf("failed to write search results: %w", err)
	}
	return resultWriter.GetResultDir(), nil
}

func newTargetKubeconfig(kubeconfig *starlarkstruct.Struct, clusterCtxName string) (*k8s.Client, error) {
	kcpConfig, err := kubeconfig.Attr("extra_kubeconfig")
	if err != nil {
		return nil, err
	}
	kcpConfigStr := kcpConfig.(starlark.String)
	kcpApiConfig, err := clientcmd.Load([]byte(kcpConfigStr.String()))
	if err != nil {
		return nil, fmt.Errorf("could not load the kubeconfig: %w", err)
	}
	restConfig, err := clientcmd.NewNonInteractiveClientConfig(*kcpApiConfig, clusterCtxName, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could build restConfig from kubeconfig: %w", err)
	}
	kcpClient, err := k8s.NewFromRestConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("could not initialize search client: %w", err)
	}
	return kcpClient, nil
}

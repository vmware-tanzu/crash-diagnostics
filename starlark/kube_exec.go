// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// KubeExecFn is a starlark built-in for executing command in target K8s pods
func KubeExecFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var namespace, pod, container, workdir, outputfile string
	var timeout int
	var command *starlark.List
	var kubeConfig *starlarkstruct.Struct

	if err := starlark.UnpackArgs(
		identifiers.kubeExec, args, kwargs,
		"namespace?", &namespace,
		"pod", &pod,
		"container?", &container,
		"cmd", &command,
		"workdir?", &workdir,
		"output_file?", &outputfile,
		"kube_config?", &kubeConfig,
		"timeout_in_seconds?", &timeout,
	); err != nil {
		return starlark.None, fmt.Errorf("failed to read args: %w", err)
	}

	if namespace == "" {
		namespace = "default"
	}
	if timeout == 0 {
		//Default timeout if not specified is 2 Minutes
		timeout = 120
	}

	if len(workdir) == 0 {
		//Defaults to crashd_config.workdir or /tmp/crashd
		if dir, err := getWorkdirFromThread(thread); err == nil {
			workdir = dir
		}
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
		return starlark.None, fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	clusterCtxName := getKubeConfigContextNameFromStruct(kubeConfig)

	execOpts := k8s.ExecOptions{
		Namespace:     namespace,
		Podname:       pod,
		ContainerName: container,
		Command:       toSlice(command),
		Timeout:       time.Duration(timeout) * time.Second,
	}
	executor, err := k8s.NewExecutor(path, clusterCtxName, execOpts)
	if err != nil {
		return starlark.None, fmt.Errorf("could not initialize search client: %w", err)
	}

	outputFilePath := filepath.Join(trimQuotes(workdir), outputfile)
	if outputfile == "" {
		outputFilePath = filepath.Join(trimQuotes(workdir), pod+".out")
	}
	err = executor.ExecCommand(ctx, outputFilePath, execOpts)

	return starlarkstruct.FromStringDict(
		starlark.String(identifiers.kubeCapture),
		starlark.StringDict{
			"file": starlark.String(outputFilePath),
			"error": func() starlark.String {
				if err != nil {
					return starlark.String(err.Error())
				}
				return ""
			}(),
		}), nil
}

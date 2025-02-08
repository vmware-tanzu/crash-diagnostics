// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// Executor is a struct that facilitates the execution of commands in Kubernetes pods.
// It uses the SPDYExecutor to stream command
type Executor struct {
	Executor remotecommand.Executor
}

type ExecOptions struct {
	Namespace     string
	Command       []string
	Podname       string
	ContainerName string
	Config        *Config
	Timeout       time.Duration
}

func NewExecutor(kubeconfig string, clusterCtxName string, opts ExecOptions) (*Executor, error) {
	restCfg, err := restConfig(kubeconfig, clusterCtxName)
	if err != nil {
		return nil, err
	}
	setCoreDefaultConfig(restCfg)
	restc, err := rest.RESTClientFor(restCfg)
	if err != nil {
		return nil, err
	}

	request := restc.Post().
		Namespace(opts.Namespace).
		Resource("pods").
		Name(opts.Podname).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: opts.ContainerName,
			Command:   opts.Command,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	executor, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	if err != nil {
		return nil, err

	}
	return &Executor{Executor: executor}, nil
}

// makeRESTConfig creates a new *rest.Config with a k8s context name if one is provided.
func restConfig(fileName, contextName string) (*rest.Config, error) {
	if fileName == "" {
		return nil, errors.New("kubeconfig file path required")
	}

	if contextName != "" {
		// create the config object from k8s config path and context
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: fileName},
			&clientcmd.ConfigOverrides{
				CurrentContext: contextName,
			}).ClientConfig()
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: fileName},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

// ExecCommand executes a command inside a specified Kubernetes pod using the SPDYExecutor.
func (k8sc *Executor) ExecCommand(ctx context.Context, outputFilePath string, execOptions ExecOptions) error {
	ctx, cancel := context.WithTimeout(ctx, execOptions.Timeout)
	defer cancel()

	file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Execute the command and stream the stdout and stderr to the file. Some commands are using stderr.
	err = k8sc.Executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: file,
		Stderr: file,
	})
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("command execution timed out. command:%s", execOptions.Command)
		}
		return err
	}

	return nil
}

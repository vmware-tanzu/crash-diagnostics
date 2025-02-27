// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// KubePortForwardrFn is a built-in starlark function that collects compute resources from a k8s cluster
// Starlark format: kube_port_forward_config(service="bar", target_port=664, [namespace="foo"])
func KubePortForwardrFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	var namespace, service string
	var targetPort int

	if err := starlark.UnpackArgs(
		identifiers.kubePortForwardConfig, args, kwargs,
		"namespace?", &namespace,
		"service", &service,
		"target_port", &targetPort,
	); err != nil {
		return starlark.None, fmt.Errorf("failed to read args: %w", err)
	}

	ctx, ok := thread.Local(identifiers.scriptCtx).(context.Context)
	if !ok || ctx == nil {
		return starlark.None, fmt.Errorf("script context not found")
	}

	kubeConfig := thread.Local(identifiers.kubeCfg).(*starlarkstruct.Struct)
	kubeConfigPath, err := getKubeConfigPathFromStruct(kubeConfig)
	if err != nil {
		return starlark.None, fmt.Errorf("failed to kubeconfig: %w", err)
	}

	client, err := k8s.New(kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("could not initialize search client: %w", err)
	}
	svc, err := client.Typed.CoreV1().Services(namespace).Get(ctx, service, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get service: %w", err)
	}

	selector := labels.SelectorFromSet(svc.Spec.Selector)

	pods, err := client.Typed.CoreV1().Pods(svc.Namespace).List(ctx, v1.ListOptions{LabelSelector: selector.String()})
	if err != nil || len(pods.Items) == 0 {
		return nil, fmt.Errorf("could not list pods: %w", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomPort int
	var found bool

	for !found {
		// get a random port between 49152-65535
		randomPort = r.Intn(16383) + 49152
		address := fmt.Sprintf(":%d", randomPort)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			continue
		}
		defer listener.Close()
		found = true
	}

	tunnelConfigDict := starlark.StringDict{
		"namespace":   starlark.String(namespace),
		"pod_name":    starlark.String(pods.Items[0].Name),
		"target_port": starlark.MakeInt(targetPort),
		"local_port":  starlark.MakeInt(randomPort),
	}

	return starlarkstruct.FromStringDict(starlark.String(identifiers.kubePortForwardConfig), tunnelConfigDict), nil
}

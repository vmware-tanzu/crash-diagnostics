// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// KubeNodesProviderFn is a built-in starlark function that collects compute resources from a k8s cluster
// Starlark format: kube_port_forward_config([namespace="foo", serviceName="bar", target_port=6643)
func KubePortForwardrFn(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	var namespace, serviceName, targetPort *starlark.String

	if err := starlark.UnpackArgs(
		identifiers.kubeNodesProvider, args, kwargs,
		"namespace?", &namespace,
		"serviceName?", &serviceName,
		"target_port", &targetPort,
	); err != nil {
		return starlark.None, errors.Wrap(err, "failed to read args")
	}

	ctx, ok := thread.Local(identifiers.scriptCtx).(context.Context)
	if !ok || ctx == nil {
		return starlark.None, fmt.Errorf("script context not found")
	}

	kubeConfig := thread.Local(identifiers.kubeCfg).(*starlarkstruct.Struct)
	kubeConfigPath, err := getKubeConfigPathFromStruct(kubeConfig)
	if err != nil {
		return starlark.None, errors.Wrap(err, "failed to kubeconfig")
	}

	client, err := k8s.New(kubeConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize search client")
	}
	service, err := client.Typed.CoreV1().Services(namespace.String()).Get(ctx, serviceName.String(), v1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "could not get service")
	}

	selector := labels.SelectorFromSet(service.Spec.Selector)

	pods, err := client.Typed.CoreV1().Pods(service.Namespace).List(ctx, v1.ListOptions{LabelSelector: selector.String()})
	if err != nil || len(pods.Items) == 0 {
		return nil, errors.Wrap(err, "could not list pods")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomPort int
	var found bool

	for !found {
		randomPort = r.Intn(10000) + 1000
		address := fmt.Sprintf(":%d", randomPort)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			continue
		}
		defer listener.Close()
		found = true
	}

	tunnelConfigDict := starlark.StringDict{
		"namespace":   starlark.String(identifiers.kubePortForwardConfig),
		"pod_name":    starlark.String(pods.Items[0].Name),
		"target_port": targetPort,
		"local_port":  starlark.MakeInt(randomPort),
	}

	return starlarkstruct.FromStringDict(starlark.String(identifiers.kubePortForwardConfig), tunnelConfigDict), nil
}

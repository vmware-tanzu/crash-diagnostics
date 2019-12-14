// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	LegacyGroupName = "core"
)

// Client prepares and exposes a dynamic, discovery, and Rest clients
type Client struct {
	Client      dynamic.Interface
	Disco       discovery.DiscoveryInterface
	CoreRest    rest.Interface
	JsonPrinter printers.JSONPrinter
}

// New returns a *Client
func New(kubeconfig string) (*Client, error) {
	// creating cfg for each client type because each
	// setup its own cfg default which may not be compatible
	dynCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := dynamic.NewForConfig(dynCfg)
	if err != nil {
		return nil, err
	}

	discoCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	disco, err := discovery.NewDiscoveryClientForConfig(discoCfg)
	if err != nil {
		return nil, err
	}

	restCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	setCoreDefaultConfig(restCfg)
	restc, err := rest.RESTClientFor(restCfg)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client, Disco: disco, CoreRest: restc}, nil
}

func setCoreDefaultConfig(config *rest.Config) {
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.APIPath = "/api"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
}

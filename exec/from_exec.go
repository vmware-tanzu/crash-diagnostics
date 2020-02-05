// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func exeFrom(k8s *k8s.Client, src *script.Script) (*script.FromCommand, []*script.Machine, error) {
	fromCmds, ok := src.Preambles[script.CmdFrom]
	if !ok {
		return nil, nil, fmt.Errorf("%s not defined", script.CmdFrom)
	}
	if len(fromCmds) < 1 {
		return nil, nil, fmt.Errorf("script missing valid %s", script.CmdFrom)
	}

	fromCmd, ok := fromCmds[0].(*script.FromCommand)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected type %T for %s", fromCmd, script.CmdFrom)
	}

	var machines []*script.Machine
	// retrieve from host list
	logrus.Debugf("Building machine list: 'FROM hosts:%s'", fromCmd.Hosts())
	for _, host := range fromCmd.Hosts() {
		var addr, port, name string
		parts := strings.Split(host, ":")
		if len(parts) > 1 {
			addr = parts[0]
			port = parts[1]
			name = host
		} else {
			addr = parts[0]
			port = fromCmd.Port()
			name = host
		}
		machines = append(machines, script.NewMachine(addr, port, name))
	}

	// continue on only with a valid K8s client
	if k8s == nil {
		return fromCmd, machines, nil
	}

	logrus.Debugf("Building machine list: 'FROM nodes:%s'", fromCmd.Nodes())
	fromNodes := fromCmd.Nodes()
	var allNodes []*coreV1.Node

	nodeStr := strings.Join(fromNodes, " ")
	if len(fromNodes) == 1 && fromNodes[0] == "all" {
		nodeStr = ""
	}

	nodes, err := getNodes(k8s, nodeStr, fromCmd.Labels())
	if err != nil {
		return fromCmd, machines, err
	}
	allNodes = append(allNodes, nodes...)

	for _, node := range allNodes {
		ip := getNodeInternalIP(node)
		port := fromCmd.Port()
		name := getNodeHostname(node)
		machine := script.NewMachine(ip, port, name)
		machines = append(machines, machine)
	}

	logrus.Debugf("Created %d machines", len(machines))

	return fromCmd, machines, nil
}

func getNodes(k8sc *k8s.Client, names, labels string) ([]*coreV1.Node, error) {
	objs, err := k8sc.Search(
		"core",  // group
		"nodes", // kind
		"",      // namespaces
		"",      // version
		names,
		labels,
		"", // containers
	)
	if err != nil {
		return nil, err
	}

	// collate
	var nodes []*coreV1.Node
	for _, obj := range objs {
		unstructList, ok := obj.(*unstructured.UnstructuredList)
		if !ok {
			return nil, fmt.Errorf("unexpected type for NodeList: %T", obj)
		}
		for _, item := range unstructList.Items {
			node := new(coreV1.Node)
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &node); err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func getNodeInternalIP(node *coreV1.Node) (ipAddr string) {
	for _, addr := range node.Status.Addresses {
		if addr.Type == "InternalIP" {
			ipAddr = addr.Address
			return
		}
	}
	return
}

func getNodeHostname(node *coreV1.Node) (hostname string) {
	for _, addr := range node.Status.Addresses {
		if addr.Type == "Hostname" {
			hostname = addr.Address
			return
		}
	}
	return
}

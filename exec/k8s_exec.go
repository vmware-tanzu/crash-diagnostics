// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func exeClusterInfo(src *script.Script, path string) {
	cfgs, ok := src.Preambles[script.CmdKubeConfig]
	if !ok {
		logrus.Warn("Skipping cluster-info, KUBECONFIG not provided")
		return
	}
	cfgCmd := cfgs[0].(*script.KubeConfigCommand)
	if _, err := os.Stat(cfgCmd.Config()); err != nil {
		logrus.Warnf("Skipping cluster-info, unable to load KUBECONFIG %s: %s", cfgCmd.Config(), err)
		return
	}

	logrus.Debugf("Using KUBECONFIG %s", cfgCmd.Config())
	k8sClient, err := k8s.GetClient(cfgCmd.Config())
	if err != nil {
		logrus.Errorf("Skipping cluster-info, failed to create Kubernetes API server client: %s", err)
		return
	}

	if err := k8s.DumpClusterInfo(k8sClient, path); err != nil {
		logrus.Errorf("Failed to retrieve cluster information: %s", err)
		return
	}
}

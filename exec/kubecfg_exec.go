// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func exeKubeConfig(src *script.Script) (*k8s.Client, error) {
	cfgs, ok := src.Preambles[script.CmdKubeConfig]
	if !ok {
		return nil, fmt.Errorf("KUBECONFIG not found in script")
	}
	cfgCmd := cfgs[0].(*script.KubeConfigCommand)
	if _, err := os.Stat(cfgCmd.Path()); err != nil {
		return nil, fmt.Errorf("path stat for %s: %s", cfgCmd.Path(), err)
	}

	logrus.Debugf("KUBECONFIG: path: %s", cfgCmd.Path())
	k8sClient, err := k8s.New(cfgCmd.Path())
	if err != nil {
		return nil, err
	}

	return k8sClient, nil
}

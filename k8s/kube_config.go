// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/vladimirvivien/gexe"
)

// FetchWorkloadConfig...
func FetchWorkloadConfig(clusterName, clusterNamespace, mgmtKubeConfigPath string) (string, error) {
	var filePath string
	cmdStr := fmt.Sprintf(`kubectl get secrets/%s-kubeconfig --template '{{.data.value}}' --namespace=%s --kubeconfig %s`, clusterName, clusterNamespace, mgmtKubeConfigPath)
	p := gexe.RunProc(cmdStr)
	if p.Err() != nil {
		return filePath, fmt.Errorf("kubectl get secrets failed: %s: %s", p.Err(), p.Result())
	}

	f, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("%s-workload-config", clusterName))
	if err != nil {
		return filePath, fmt.Errorf("Cannot create temporary file: %w", err)
	}
	filePath = f.Name()
	defer f.Close()

	base64Dec := base64.NewDecoder(base64.StdEncoding, p.Out())
	if _, err := io.Copy(f, base64Dec); err != nil {
		return filePath, fmt.Errorf("error decoding workload kubeconfig: %w", err)
	}
	return filePath, nil
}

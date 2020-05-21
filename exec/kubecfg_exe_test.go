// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/parser"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestExecKUBECONFIGFunc(t *testing.T) {
	clusterName := "crashd-test-kubecfg"
	k8sconfig := fmt.Sprintf("/tmp/%s", clusterName)

	// create kind cluster
	kind := testcrashd.NewKindCluster("../testing/kind-cluster-docker.yaml", clusterName)
	if err := kind.Create(); err != nil {
		t.Fatal(err)
	}
	defer kind.Destroy()

	if err := kind.MakeKubeConfigFile(k8sconfig); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(k8sconfig)

	tests := []struct {
		name   string
		script func() *script.Script
		exec   func(*script.Script)
	}{
		{
			name: "KUBECONFIG with path OK",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
				`, k8sconfig)
				script, _ := parser.Parse(strings.NewReader(src))
				return script
			},
			exec: func(src *script.Script) {
				k8sc, err := exeKubeConfig(src)
				if err != nil {
					t.Fatal(err)
				}
				if k8sc == nil {
					t.Fatalf("Unexpected nil k8sc.Client")
				}
			},
		},
		{
			name: "KUBECONFIG with bad path",
			script: func() *script.Script {
				src := fmt.Sprintf(`KUBECONFIG bad-path`)
				script, _ := parser.Parse(strings.NewReader(src))
				return script
			},
			exec: func(src *script.Script) {
				k8sc, err := exeKubeConfig(src)
				if err == nil {
					t.Fatal("Expecting exeKubeConfig to fail, but didnt")
				}
				if k8sc != nil {
					t.Fatalf("Expected nil k8sc.Client, but it's not")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.exec(test.script())
		})
	}
}

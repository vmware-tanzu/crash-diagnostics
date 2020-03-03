// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestExecKUBEGETFunc(t *testing.T) {
	clusterName := "crashd-test-kubeget"
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

	// important, wait for at least 1 pod to be deployed
	waitTime := time.Second * 11
	logrus.Infof("Sleeping %v ... waiting for pods", waitTime)
	time.Sleep(waitTime)

	tests := []struct {
		name   string
		script func() *script.Script
		exec   func(*k8s.Client, *script.Script)
	}{
		{
			name: "KUBEGET pods",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					KUBEGET objects groups:"core" kinds:"pods" namespaces:"kube-system"
				`, k8sconfig)
				script, err := script.Parse(strings.NewReader(src))
				if err != nil {
					t.Fatal(err)
				}
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) {
				if k8sc == nil {
					t.Log("k8s.Client == nil, skipping test")
					return
				}
				cmd0, ok := src.Actions[0].(*script.KubeGetCommand)
				if !ok {
					t.Fatalf("Unexpected script action type for %T", cmd0)
				}
				objects, err := exeKubeGet(k8sc, cmd0)
				if err != nil {
					t.Fatal(err)
				}
				if len(objects) == 0 {
					t.Fatal("exeKubeGet returns 0 objects")
				}
			},
		},
		{
			name: "KUBEGET pods with labels",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					KUBEGET objects groups:"core" kinds:"pods" namespaces:"kube-system" labels:"component=kube-apiserver"
				`, k8sconfig)
				script, err := script.Parse(strings.NewReader(src))
				if err != nil {
					t.Fatal(err)
				}
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) {
				if k8sc == nil {
					t.Log("k8s.Client == nil, skipping test")
					return
				}
				cmd0, ok := src.Actions[0].(*script.KubeGetCommand)
				if !ok {
					t.Fatalf("Unexpected script action type for %T", cmd0)
				}
				objects, err := exeKubeGet(k8sc, cmd0)
				if err != nil {
					t.Fatal(err)
				}
				if len(objects) != 1 {
					t.Fatalf("exeKubeGet got unexpected number of objects %d", len(objects))
				}
			},
		},
		{
			name: "KUBEGET pod logs",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					KUBEGET logs groups:"core" kinds:"pods" namespaces:"kube-system" labels:"component=kube-apiserver"
				`, k8sconfig)
				script, err := script.Parse(strings.NewReader(src))
				if err != nil {
					t.Fatal(err)
				}
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) {
				if k8sc == nil {
					t.Log("k8s.Client == nil, skipping test")
					return
				}
				cmd0, ok := src.Actions[0].(*script.KubeGetCommand)
				if !ok {
					t.Fatalf("Unexpected script action type for %T", cmd0)
				}
				objects, err := exeKubeGet(k8sc, cmd0)
				if err != nil {
					t.Fatal(err)
				}
				if len(objects) != 1 {
					t.Fatalf("exeKubeGet got unexpected number of objects %d", len(objects))
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src := test.script()
			k8sc, err := exeKubeConfig(src)
			if err != nil {
				t.Log(err)
			}
			test.exec(k8sc, src)
		})
	}
}

func TestExecKUBEGET(t *testing.T) {
	clusterName := "crashd-test-kubeget"
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

	// important, wait for at least 1 pod to be deployed
	waitTime := time.Second * 11
	logrus.Infof("Sleeping %v ... waiting for pods", waitTime)
	time.Sleep(waitTime)

	tests := []execTest{
		{
			name: "KUBEGET",
			source: func() string {
				return `
				FROM local
				KUBECONFIG $HOME/.kube/kind-config-kind
				KUBEGET objects namespaces:"kube-system"
				`
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

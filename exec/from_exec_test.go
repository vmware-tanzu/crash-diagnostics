// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/k8s"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestExecFROMFunc(t *testing.T) {
	clusterName := "crashd-test-from"
	kindNodeName := fmt.Sprintf("%s-control-plane", clusterName)
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

	// tests
	tests := []struct {
		name   string
		script func() *script.Script
		exec   func(*k8s.Client, *script.Script) error
	}{
		{
			name: "FROM with host:port",
			script: func() *script.Script {
				script, _ := script.Parse(strings.NewReader("FROM 1.1.1.1:4444"))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				fromCmd, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != len(fromCmd.Hosts()) {
					return fmt.Errorf("FROM: expecting %d machines got %d", len(fromCmd.Hosts()), len(machines))
				}
				machine := machines[0]
				if machine.Host() != "1.1.1.1" {
					return fmt.Errorf("FROM machine has unexpected host %s", machine.Host())
				}
				if machine.Port() != "4444" {
					return fmt.Errorf("FROM machine has unexpected port %s", machine.Port())
				}

				return nil
			},
		},
		{
			name: "FROM with host default port",
			script: func() *script.Script {
				script, _ := script.Parse(strings.NewReader("FROM 1.1.1.1"))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				fromCmd, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != len(fromCmd.Hosts()) {
					return fmt.Errorf("FROM: expecting %d machines got %d", len(fromCmd.Hosts()), len(machines))
				}
				machine := machines[0]
				if machine.Host() != "1.1.1.1" {
					return fmt.Errorf("FROM machine has unexpected host %s", machine.Host())
				}
				if machine.Port() != "22" {
					return fmt.Errorf("FROM machine has unexpected port %s", machine.Port())
				}

				return nil
			},
		},
		{
			name: "FROM with host:port and global port",
			script: func() *script.Script {
				script, _ := script.Parse(strings.NewReader(`FROM hosts:"1.1.1.1 10.10.10.10:2222" port:2121`))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				fromCmd, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != len(fromCmd.Hosts()) {
					return fmt.Errorf("FROM: expecting %d machines got %d", len(fromCmd.Hosts()), len(machines))
				}
				m0 := machines[0]
				m1 := machines[1]
				if m0.Host() != "1.1.1.1" || m0.Port() != "2121" {
					return fmt.Errorf("FROM machine0 has unexpected host:port %s:%s", m0.Host(), m0.Port())
				}
				if m1.Host() != "10.10.10.10" || m1.Port() != "2222" {
					return fmt.Errorf("FROM machine1 has unexpected host:port %s:%s", m1.Host(), m1.Port())
				}

				return nil
			},
		},
		{
			name: "FROM with all nodes",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					FROM nodes:'all'
				`, k8sconfig)
				script, _ := script.Parse(strings.NewReader(src))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				fromCmd, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != 1 {
					return fmt.Errorf("FROM %#v: expecting 1 machine got %d", fromCmd.Args(), len(machines))
				}

				machine := machines[0]
				t.Logf("Machine found: %#v", machine)

				if machine.Port() != fromCmd.Port() {
					return fmt.Errorf("FROM machine has unexpected port %s", machine.Port())
				}

				if machine.Name() != kindNodeName {
					return fmt.Errorf("FROM machine has unexpected node name %s", machine.Name())
				}

				return nil
			},
		},
		{
			name: "FROM with specific nodes",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					FROM nodes:'%s'
				`, k8sconfig, kindNodeName)
				script, _ := script.Parse(strings.NewReader(src))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				fromCmd, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != 1 {
					return fmt.Errorf("FROM %#v: expecting 1 machine got %d", fromCmd.Args(), len(machines))
				}

				machine := machines[0]

				if machine.Name() != kindNodeName {
					return fmt.Errorf("FROM machine has unexpected node name %s", machine.Name())
				}

				return nil
			},
		},
		{
			name: "FROM with bad node name",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					FROM nodes:'bad-node-name'
				`, k8sconfig)
				script, _ := script.Parse(strings.NewReader(src))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				_, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != 0 {
					return fmt.Errorf("FROM: expecting 0 machine, got %d", len(machines))
				}
				return nil
			},
		},
		{
			name: "FROM with node labels",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					FROM nodes:'all' labels:'kubernetes.io/hostname=%s'
				`, k8sconfig, kindNodeName)
				script, _ := script.Parse(strings.NewReader(src))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				fromCmd, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != 1 {
					return fmt.Errorf("FROM %#v: expecting 1 machine got %d", fromCmd.Args(), len(machines))
				}

				machine := machines[0]

				if machine.Name() != kindNodeName {
					return fmt.Errorf("FROM machine has unexpected node name %s", machine.Name())
				}

				return nil
			},
		},
		{
			name: "FROM with bad node labels",
			script: func() *script.Script {
				src := fmt.Sprintf(`
					KUBECONFIG %s
					FROM nodes:'all' labels:'foo/bar=mycluster-control-plane'
				`, k8sconfig)
				script, _ := script.Parse(strings.NewReader(src))
				return script
			},
			exec: func(k8sc *k8s.Client, src *script.Script) error {
				_, machines, err := exeFrom(k8sc, src)
				if err != nil {
					return err
				}
				if len(machines) != 0 {
					return fmt.Errorf("FROM: expecting 0 machine got %d", len(machines))
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src := test.script()
			k8sc, err := exeKubeConfig(src)
			if err != nil {
				t.Logf("Failed to get KubeConfig: %s", err)
			}
			if err := test.exec(k8sc, src); err != nil {
				t.Error(err)
			}
		})
	}

}

func TestExecFROM(t *testing.T) {
	tests := []execTest{
		{
			name: "FROM with multiple addresses",
			source: func() string {
				return `
				ENV host=local
				FROM '$host'
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

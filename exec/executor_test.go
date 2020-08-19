// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

var (
	support *testcrashd.TestSupport
)

func TestMain(m *testing.M) {
	test, err := testcrashd.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	support = test

	if err := support.SetupSSHServer(); err != nil {
		logrus.Fatal(err)
	}

	if err := support.SetupKindCluster(); err != nil {
		logrus.Fatal(err)
	}

	_, err = support.SetupKindKubeConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	result := m.Run()

	if err := support.TearDown(); err != nil {
		logrus.Fatal(err)
	}

	os.Exit(result)
}

func TestKindScript(t *testing.T) {
	tests := []struct {
		name       string
		scriptPath string
		args       ArgMap
	}{
		{
			name:       "api objects",
			scriptPath: "../examples/kind-api-objects.crsh",
			args:       ArgMap{"kubecfg": support.KindKubeConfigFile()},
		},
		{
			name:       "pod logs",
			scriptPath: "../examples/pod-logs.crsh",
			args:       ArgMap{"kubecfg": support.KindKubeConfigFile()},
		},
		{
			name:       "script with args",
			scriptPath: "../examples/script-args.crsh",
			args: ArgMap{
				"workdir": "/tmp/crashargs",
				"kubecfg": support.KindKubeConfigFile(),
				"output":  "/tmp/craslogs.tar.gz",
			},
		},
		{
			name:       "host-list provider",
			scriptPath: "../examples/host-list-provider.crsh",
			args: ArgMap{
				"kubecfg":     support.KindKubeConfigFile(),
				"ssh_pk_path": support.PrivateKeyPath(),
				"ssh_port":    support.PortValue(),
			},
		},
		//{
		//	name:       "kube-nodes provider",
		//	scriptPath: "../examples/kube-nodes-provider.crsh",
		//	args: ArgMap{
		//		"kubecfg":  getTestKubeConf(),
		//		"ssh_port": testSSHPort,
		//		"username": getUsername(),
		//		"key_path": getPrivateKey(),
		//	},
		//},
		{
			name:       "kind-capi-bootstrap",
			scriptPath: "../examples/kind-capi-bootstrap.crsh",
			args:       ArgMap{"cluster_name": support.ResourceName()},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.Open(test.scriptPath)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			if err := ExecuteFile(file, test.args); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name   string
		script string
		exec   func(t *testing.T, script string)
	}{
		{
			name:   "run_local",
			script: `result = run_local("echo 'Hello World!'")`,
			exec: func(t *testing.T, script string) {
				if err := Execute("run_local", strings.NewReader(script), ArgMap{}); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.exec(t, test.script)
		})
	}
}

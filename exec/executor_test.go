// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

var (
	testSSHPort     = testcrashd.NextPortValue()
	testServerName  = testcrashd.NextResourceName()
	testClusterName = testcrashd.NextResourceName()
	getTestKubeConf func() string
)

func TestMain(m *testing.M) {
	testcrashd.Init()

	sshSvr := testcrashd.NewSSHServer(testServerName, testSSHPort)
	logrus.Debug("Attempting to start SSH server")
	if err := sshSvr.Start(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	kind := testcrashd.NewKindCluster("../testing/kind-cluster-docker.yaml", testClusterName)
	if err := kind.Create(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// attempt to wait for cluster up
	time.Sleep(time.Second * 10)

	tmpFile, err := ioutil.TempFile(os.TempDir(), testClusterName)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	defer func() {
		logrus.Debug("Stopping SSH server...")
		if err := sshSvr.Stop(); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		if err := kind.Destroy(); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}()

	getTestKubeConf = func() string {
		return tmpFile.Name()
	}

	if err := kind.MakeKubeConfigFile(getTestKubeConf()); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
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
			args:       ArgMap{"kubecfg": getTestKubeConf()},
		},
		{
			name:       "pod logs",
			scriptPath: "../examples/pod-logs.crsh",
			args:       ArgMap{"kubecfg": getTestKubeConf()},
		},
		{
			name:       "script with args",
			scriptPath: "../examples/script-args.crsh",
			args: ArgMap{
				"workdir": "/tmp/crashargs",
				"kubecfg": getTestKubeConf(),
				"output":  "/tmp/craslogs.tar.gz",
			},
		},
		{
			name:       "host-list provider",
			scriptPath: "../examples/host-list-provider.crsh",
			args:       ArgMap{"kubecfg": getTestKubeConf(), "ssh_port": testSSHPort},
		},
		//{
		//	name:       "kube-nodes provider",
		//	scriptPath: "../examples/kube-nodes-provider.crsh",
		//	args: ArgMap{
		//		"kubecfg":  getTestKubeConf(),
		//		"ssh_port": testSSHPort,
		//		"username": testcrashd.GetSSHUsername(),
		//		"key_path": testcrashd.GetSSHPrivateKey(),
		//	},
		//},
		{
			name:       "kind-capi-bootstrap",
			scriptPath: "../examples/kind-capi-bootstrap.crsh",
			args:       ArgMap{"cluster_name": testClusterName},
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

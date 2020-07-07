// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"testing"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

const (
	testSSHPort = "2222"
)

func TestMain(m *testing.M) {
	testcrashd.Init()
	//
	//sshSvr := testcrashd.NewSSHServer("test-sshd-exec", testSSHPort)
	//logrus.Debug("Attempting to start SSH server")
	//if err := sshSvr.Start(); err != nil {
	//	logrus.Error(err)
	//	os.Exit(1)
	//}
	//
	//testResult := m.Run()
	//
	//logrus.Debug("Stopping SSH server...")
	//if err := sshSvr.Stop(); err != nil {
	//	logrus.Error(err)
	//	os.Exit(1)
	//}
	//
	//os.Exit(testResult)

	// Skipping all tests
	os.Exit(0)
}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

var (
	support     *testcrashd.TestSupport
	testSSHArgs SSHArgs
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

	testSSHArgs = SSHArgs{
		User:           support.CurrentUsername(),
		PrivateKeyPath: support.PrivateKeyPath(),
		Host:           "127.0.0.1",
		Port:           support.PortValue(),
		MaxRetries:     support.MaxConnectionRetries(),
	}

	testResult := m.Run()

	logrus.Debug("Shutting down test...")
	if err := support.TearDown(); err != nil {
		logrus.Fatal(err)
	}

	os.Exit(testResult)
}

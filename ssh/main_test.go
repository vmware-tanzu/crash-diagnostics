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
	testSSHPort    = testcrashd.NextPortValue()
	testMaxRetries = 30
)

func TestMain(m *testing.M) {
	testcrashd.Init()

	sshSvr := testcrashd.NewSSHServer(testcrashd.NextResourceName(), testSSHPort)
	logrus.Debug("Attempting to start SSH server")
	if err := sshSvr.Start(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	testResult := m.Run()

	logrus.Debug("Stopping SSH server...")
	if err := sshSvr.Stop(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	os.Exit(testResult)
}

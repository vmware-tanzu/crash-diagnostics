// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
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
		logrus.Fatal("failed to initialize test support:", err)
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

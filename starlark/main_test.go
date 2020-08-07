// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

var (
	testSupport *testcrashd.TestSupport
)

func TestMain(m *testing.M) {
	test, err := testcrashd.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	testSupport = test

	if err := testSupport.SetupSSHServer(); err != nil {
		logrus.Fatal(err)
	}

	if err := testSupport.SetupKindCluster(); err != nil {
		logrus.Fatal(err)
	}

	// precaution
	if testSupport == nil {
		logrus.Fatal("failed to setup test support")
	}

	result := m.Run()

	if err := testSupport.TearDown(); err != nil {
		logrus.Fatal(err)
	}

	os.Exit(result)
}

func makeTestSSHConfig(pkPath, port, username string) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
		identifiers.username:       starlark.String(username),
		identifiers.port:           starlark.String(port),
		identifiers.privateKeyPath: starlark.String(pkPath),
		identifiers.maxRetries:     starlark.String(fmt.Sprintf("%d", testSupport.MaxConnectionRetries())),
	})
}

func makeTestSSHHostResource(addr string, sshCfg *starlarkstruct.Struct) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		starlark.StringDict{
			"kind":       starlark.String(identifiers.hostResource),
			"provider":   starlark.String(identifiers.hostListProvider),
			"host":       starlark.String(addr),
			"transport":  starlark.String("ssh"),
			"ssh_config": sshCfg,
		},
	)
}

func newTestThreadLocal(t *testing.T) *starlark.Thread {
	thread := &starlark.Thread{Name: "test-crashd"}
	if err := setupLocalDefaults(thread); err != nil {
		t.Fatalf("failed to setup new thread local: %s", err)
	}
	return thread
}

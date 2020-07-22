// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"os"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestMain(m *testing.M) {
	testcrashd.Init()
	os.Exit(m.Run())
}

func makeTestSSHConfig(pkPath, port string) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
		identifiers.username:       starlark.String(testcrashd.GetSSHUsername()),
		identifiers.port:           starlark.String(port),
		identifiers.privateKeyPath: starlark.String(pkPath),
		identifiers.maxRetries:     starlark.String(defaults.connRetries),
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

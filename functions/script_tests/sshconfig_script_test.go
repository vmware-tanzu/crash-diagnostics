// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script_tests

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

func TestSSHConfScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "make ssh config",
			script: `
result = make_ssh_config(username="foo", port="44", private_key_path="./ssh/path", max_retries=32)
`,
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}

				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("script_conf() should be assigned to a variable for test")
				}
				var result sshconf.Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}
				conf := result.Config
				if conf.Username != "foo" {
					t.Errorf("unexpected username value: %s", conf.Username)
				}
				if conf.Port != "44" {
					t.Errorf("unexpected port value: %s", conf.Port)
				}
				if conf.PrivateKeyPath != "./ssh/path" {
					t.Errorf("unexpected pk path value: %s", conf.PrivateKeyPath)
				}
				if conf.MaxRetries != 32 {
					t.Errorf("unexpected max retries value: %d", conf.MaxRetries)
				}
			},
		},
		{
			name: "ssh config alias",
			script: `
result = ssh_config(username="foo", port="44", private_key_path="./ssh/path", max_retries=32)
`,
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}

				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("script_conf() should be assigned to a variable for test")
				}
				var conf sshconf.Config
				if err := typekit.Starlark(resultVal).Go(&conf); err != nil {
					t.Fatal(err)
				}

				if conf.Username != "foo" {
					t.Errorf("unexpected username value: %s", conf.Username)
				}
				if conf.Port != "44" {
					t.Errorf("unexpected port value: %s", conf.Port)
				}
				if conf.PrivateKeyPath != "./ssh/path" {
					t.Errorf("unexpected pk path value: %s", conf.PrivateKeyPath)
				}
				if conf.MaxRetries != 32 {
					t.Errorf("unexpected max retries value: %d", conf.MaxRetries)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script)
		})
	}
}

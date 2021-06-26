// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package scriptconf

import (
	"os"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"go.starlark.net/starlark"
)

func TestScriptConfRun(t *testing.T) {
	tests := []struct {
		name   string
		params Args
		config Config
	}{
		{
			name:   "default values",
			params: Args{},
			config: Config{Workdir: DefaultWorkdir(), Gid: functions.DefaultGid(), Uid: functions.DefaultUid()},
		},
		{
			name:   "all values",
			params: Args{Workdir: "foo", Gid: "00", Uid: "01", UseSSHAgent: true, Requires: []string{"a/b"}},
			config: Config{Workdir: "foo", Gid: "00", Uid: "01", UseSSHAgent: true, Requires: []string{"a/b"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			thread := &starlark.Thread{}
			result := Run(thread, test.params)
			if result.Error != "" {
				t.Fatal(result.Error)
			}

			cfg := result.Conf
			if cfg.Workdir != test.config.Workdir {
				t.Errorf("unexpected workdir value %s", cfg.Workdir)
			}
			if err := os.RemoveAll(test.config.Workdir); err != nil {
				t.Fatal(err)
			}
			if cfg.Gid != test.config.Gid {
				t.Errorf("unexpected Gid: %s", cfg.Gid)
			}
			if cfg.Uid != test.config.Uid {
				t.Errorf("expected Uid %s, got: %s", test.config.Uid, cfg.Uid)
			}
			if cfg.UseSSHAgent != test.config.UseSSHAgent {
				t.Errorf("unexpected UseSSHAgent: %t", cfg.UseSSHAgent)
			}
			if cfg.UseSSHAgent {
				if thread.Local(sshconf.AgentIdentifier) == nil {
					t.Errorf("ssh_agent was not stored in thread_local")
				}
			}
			if len(cfg.Requires) != len(test.config.Requires) {
				t.Errorf("unexpected len(Requires) %d", len(cfg.Requires))
			}
		})
	}
}

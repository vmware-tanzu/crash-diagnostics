// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sshconf

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestConfCmd_Run(t *testing.T) {
	tests := []struct {
		name       string
		args       Args
		config     Config
		shouldFail bool
	}{
		{
			name:       "zero values",
			args:       Args{},
			shouldFail: true,
		},
		{
			name: "default values",
			args: Args{Username: "testuser"},
			config: Config{
				Username:       "testuser",
				Port:           "22",
				PrivateKeyPath: DefaultPKPath(),
				JumpUsername:   "",
				JumpHost:       "",
				MaxRetries:     0,
				ConnTimeout:    DefaultConnTimeout(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			thread := &starlark.Thread{}
			result := newCmd().Run(thread, test.args)
			if result.Error != "" && !test.shouldFail {
				t.Fatal(result.Error)
			}

			cfg := result
			if cfg.ConnTimeout != test.config.ConnTimeout {
				t.Errorf("unexpected conntimeout value %d", cfg.ConnTimeout)
			}
			if cfg.MaxRetries != test.config.MaxRetries {
				t.Errorf("unexpected max retries: %d", cfg.MaxRetries)
			}
			if cfg.JumpHost != test.config.JumpHost {
				t.Errorf("unexpected jump host: %s", cfg.JumpHost)
			}
			if cfg.JumpUsername != test.config.JumpUsername {
				t.Errorf("unexpected jump username: %s", cfg.JumpUsername)
			}
			if cfg.Port != test.config.Port {
				t.Errorf("unexpected ssh port: %s", cfg.Port)
			}
		})
	}
}

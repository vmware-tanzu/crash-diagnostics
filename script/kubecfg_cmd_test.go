// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"testing"
)

func TestCommandKUBECONFIG(t *testing.T) {
	tests := []commandTest{
		{
			name: "KUBECONFIG with single path",
			source: func() string {
				return "KUBECONFIG /a/b/c"
			},
			script: func(s *Script) error {
				cfgs := s.Preambles[CmdKubeConfig]
				if len(cfgs) != 1 {
					return fmt.Errorf("Script has unexpected number of KUBECONFIG %d", len(cfgs))
				}
				cfg, ok := cfgs[0].(*KubeConfigCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", cfgs[0])
				}
				if cfg.Config() != "/a/b/c" {
					return fmt.Errorf("KUBECONFIG has unexpected config %s", cfg.Config())
				}
				return nil
			},
		},
		{
			name: "Script with multiple KUBECONFIG",
			source: func() string {
				return "KUBECONFIG /a/b/c\nKUBECONFIG /e/f/g"
			},
			script: func(s *Script) error {
				cfgs := s.Preambles[CmdKubeConfig]
				if len(cfgs) != 1 {
					return fmt.Errorf("Script has unexpected number of KUBECONFIG %d", len(cfgs))
				}
				cfg, ok := cfgs[0].(*KubeConfigCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", cfgs[0])
				}
				if cfg.Config() != "/e/f/g" {
					return fmt.Errorf("KUBECONFIG has unexpected config %s", cfg.Config())
				}
				return nil
			},
		},
		{
			name: "KUBECONFIG with multiple paths",
			source: func() string {
				return "KUBECONFIG /a/b/c /d/e/f"
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}

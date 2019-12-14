// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"testing"
)

func TestCommandKUBEGET(t *testing.T) {
	tests := []commandTest{
		{
			name: "KUBEGET unamed what-param, no other params",
			source: func() string {
				return "KUBEGET objects"
			},
			script: func(s *Script) error {
				kgCmd, ok := s.Actions[0].(*KubeGetCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", s.Actions[0])
				}
				if kgCmd.What() != "objects" {
					return fmt.Errorf("KUBEGET unexpected what: %s", kgCmd.What())
				}
				return nil
			},
		},
		{
			name: "KUBEGET named what-param, no other params",
			source: func() string {
				return "KUBEGET what:logs"
			},
			script: func(s *Script) error {
				kgCmd, ok := s.Actions[0].(*KubeGetCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", s.Actions[0])
				}
				if kgCmd.What() != "logs" {
					return fmt.Errorf("KUBEGET unexpected what: %s", kgCmd.What())
				}
				return nil
			},
		},
		{
			name: "KUBEGET objects with other params",
			source: func() string {
				return `
				KUBEGET objects namespaces:"myns testns" groups:"v1" kinds:"pods events" versions:"1" names:"my-app" labels:"prod" containers:"webapp"`
			},
			script: func(s *Script) error {
				kgCmd := s.Actions[0].(*KubeGetCommand)
				if len(kgCmd.Args()) != 8 {
					return fmt.Errorf("KUBEGET unexpected param count: %d", len(kgCmd.Args()))
				}
				// check each param
				if kgCmd.What() != "objects" {
					return fmt.Errorf("KUBEGET unexpected what: %s", kgCmd.What())
				}
				if kgCmd.Namespaces() != "myns testns" {
					return fmt.Errorf("KUBEGET unexpected namespaces: %s", kgCmd.Namespaces())
				}
				if kgCmd.Groups() != "v1" {
					return fmt.Errorf("KUBEGET unexpected namespaces: %s", kgCmd.Namespaces())
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}

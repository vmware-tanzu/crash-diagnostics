// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCommandFROM(t *testing.T) {
	tests := []commandTest{
		{
			name: "FROM default hosts param unquoted",
			source: func() string {
				return "FROM local foo.bar:1234"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd, ok := froms[0].(*FromCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", froms[0])
				}
				if len(fromCmd.Hosts()) != 2 {
					return fmt.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Nodes()))
				}
				if len(fromCmd.Nodes()) != 0 {
					return fmt.Errorf("FROM has unexpected nodes param %v", len(fromCmd.Nodes()))
				}
				if fromCmd.Hosts()[0] != "local" && fromCmd.Hosts()[1] != "foo.basr:1234" {
					return fmt.Errorf("FROM has unexpected host address %v", fromCmd.Hosts())
				}
				// check defaults
				if fromCmd.Port() != Defaults.ServicePort {
					return fmt.Errorf("FROM has unexpected default port %s", fromCmd.Port())
				}
				if fromCmd.ConnectionRetries() != 30 {
					return fmt.Errorf("FROM has unexpected retries %d", fromCmd.ConnectionRetries())
				}
				if fromCmd.ConnectionTimeout() != time.Second*120 {
					return fmt.Errorf("FROM has unexpected retries %d", fromCmd.ConnectionRetries())
				}
				return nil
			},
		},
		{
			name: "FROM default hosts param quoted",
			source: func() string {
				return "FROM 'local foo.bar:1234'"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)
				if len(fromCmd.Hosts()) != 2 {
					return fmt.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Nodes()))
				}
				if len(fromCmd.Nodes()) != 0 {
					return fmt.Errorf("FROM has unexpected nodes param %v", fromCmd.Nodes())
				}
				if fromCmd.Hosts()[0] != "local" && fromCmd.Hosts()[1] != "foo.basr:1234" {
					return fmt.Errorf("FROM has unexpected host address %v", fromCmd.Hosts())
				}

				return nil
			},
		},
		{
			name: "FROM with nodes ports timeout",
			source: func() string {
				return "FROM nodes:'node.1 node.2 10.10.10.12' port:2222 retries:100 timeout:'5m'"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)
				if len(fromCmd.Hosts()) != 0 {
					return fmt.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Hosts()))
				}
				if len(fromCmd.Nodes()) != 3 {
					return fmt.Errorf("FROM has unexpected nodes param %v", fromCmd.Nodes())
				}
				if fromCmd.Port() != "2222" {
					return fmt.Errorf("FROM has unexpected port %s", fromCmd.Port())
				}
				if fromCmd.ConnectionRetries() != 100 {
					return fmt.Errorf("FROM has unexpected connection retries %d", fromCmd.ConnectionRetries())
				}
				if fromCmd.ConnectionTimeout() != time.Minute*5 {
					return fmt.Errorf("FROM has unexpected connection retries %d", fromCmd.ConnectionRetries())
				}
				return nil
			},
		},
		{
			name: "Multiple FROMs last-one-win",
			source: func() string {
				return `
				FROM local foo.bar:1234
				FROM nodes:'local.1:123 local.2:456'`
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)
				if len(fromCmd.Hosts()) != 0 {
					return fmt.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Hosts()))
				}
				if len(fromCmd.Nodes()) != 2 {
					return fmt.Errorf("FROM has unexpected nodes param %v", len(fromCmd.Nodes()))
				}
				if fromCmd.Nodes()[0] != "local.1:123" && fromCmd.Hosts()[1] != "local.2:456" {
					return fmt.Errorf("FROM has unexpected host address %v", fromCmd.Hosts())
				}
				return nil
			},
		},

		{
			name: "FROM with var expansion",
			source: func() string {
				os.Setenv("foohost", "foo.bar")
				os.Setenv("port", "1234")
				return "FROM hosts:${foohost} port:$port"
			},
			script: func(s *Script) error {
				froms := s.Preambles[CmdFrom]
				if len(froms) != 1 {
					return fmt.Errorf("Script has unexpected number of FROM %d", len(froms))
				}
				fromCmd := froms[0].(*FromCommand)

				if len(fromCmd.Hosts()) != 1 {
					return fmt.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Hosts()))
				}
				if fromCmd.Hosts()[0] != "foo.bar" {
					return fmt.Errorf("FROM has unexpected host value %s", fromCmd.Hosts()[0])
				}
				if fromCmd.Port() != "1234" {
					return fmt.Errorf("FROM has unexpected port value %s", fromCmd.Port())
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

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"testing"
	"time"
)

func TestCommandFROM(t *testing.T) {
	tests := []commandTest{
		{
			name: "FROM",
			command: func(t *testing.T) Command{
				cmd, err := NewFromCommand(0,"local foo.bar:1234")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				fromCmd, ok := c.(*FromCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if len(fromCmd.Hosts()) != 2 {
					t.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Nodes()))
				}
				if len(fromCmd.Nodes()) != 0 {
					t.Errorf("FROM has unexpected nodes param %v", len(fromCmd.Nodes()))
				}
				if fromCmd.Hosts()[0] != "local" && fromCmd.Hosts()[1] != "foo.basr:1234" {
					t.Errorf("FROM has unexpected host address %v", fromCmd.Hosts())
				}
				// check defaults
				if fromCmd.Port() != Defaults.ServicePort {
					t.Errorf("FROM has unexpected default port %s", fromCmd.Port())
				}
				if fromCmd.ConnectionRetries() != 30 {
					t.Errorf("FROM has unexpected retries %d", fromCmd.ConnectionRetries())
				}
				if fromCmd.ConnectionTimeout() != time.Second*120 {
					t.Errorf("FROM has unexpected retries %d", fromCmd.ConnectionRetries())
				}
				
			},
		},
		{
			name: "FROM/quoted",
			command: func(t *testing.T) Command {
				cmd, err := NewFromCommand(0, "'local foo.bar:1234'")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				fromCmd := c.(*FromCommand)
				if len(fromCmd.Hosts()) != 2 {
					t.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Nodes()))
				}
				if len(fromCmd.Nodes()) != 0 {
					t.Errorf("FROM has unexpected nodes param %v", fromCmd.Nodes())
				}
				if fromCmd.Hosts()[0] != "local" && fromCmd.Hosts()[1] != "foo.basr:1234" {
					t.Errorf("FROM has unexpected host address %v", fromCmd.Hosts())
				}

			},
		},
		{
			name: "FROM/all params",
			command: func(t *testing.T) Command {
				cmd, err := NewFromCommand(0, "nodes:'node.1 node.2 10.10.10.12' port:2222 retries:100 timeout:'5m'")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				fromCmd := c.(*FromCommand)
				if len(fromCmd.Hosts()) != 0 {
					t.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Hosts()))
				}
				if len(fromCmd.Nodes()) != 3 {
					t.Errorf("FROM has unexpected nodes param %v", fromCmd.Nodes())
				}
				if fromCmd.Port() != "2222" {
					t.Errorf("FROM has unexpected port %s", fromCmd.Port())
				}
				if fromCmd.ConnectionRetries() != 100 {
					t.Errorf("FROM has unexpected connection retries %d", fromCmd.ConnectionRetries())
				}
				if fromCmd.ConnectionTimeout() != time.Minute*5 {
					t.Errorf("FROM has unexpected connection retries %d", fromCmd.ConnectionRetries())
				}
				
			},
		},

		{
			name: "FROM/var expansion",
			command: func(t *testing.T) Command {
				cmd, err := NewFromCommand(0, "hosts:${foohost} port:$port")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				os.Setenv("foohost", "foo.bar")
				os.Setenv("port", "1234")
				
				fromCmd := c.(*FromCommand)

				if len(fromCmd.Hosts()) != 1 {
					t.Errorf("FROM has unexpected number of hosts %d", len(fromCmd.Hosts()))
				}
				if fromCmd.Hosts()[0] != "foo.bar" {
					t.Errorf("FROM has unexpected host value %s", fromCmd.Hosts()[0])
				}
				if fromCmd.Port() != "1234" {
					t.Errorf("FROM has unexpected port value %s", fromCmd.Port())
				}
				
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}

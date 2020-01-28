// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecFROMFunc(t *testing.T) {
	tests := []struct {
		name   string
		script func() *script.Script
		exec   func(*script.Script) error
	}{
		{
			name: "FROM with host:port",
			script: func() *script.Script {
				script, _ := script.Parse(strings.NewReader("FROM 1.1.1.1:4444"))
				return script
			},
			exec: func(src *script.Script) error {
				fromCmd, machines, err := exeFrom(src)
				if err != nil {
					return err
				}
				if len(machines) != len(fromCmd.Hosts()) {
					return fmt.Errorf("FROM: expecting %d machines got %d", len(fromCmd.Hosts()), len(machines))
				}
				machine := machines[0]
				if machine.Host() != "1.1.1.1" {
					return fmt.Errorf("FROM machine has unexpected host %s", machine.Host())
				}
				if machine.Port() != "4444" {
					return fmt.Errorf("FROM machine has unexpected port %s", machine.Port())
				}

				return nil
			},
		},
		{
			name: "FROM with host default port",
			script: func() *script.Script {
				script, _ := script.Parse(strings.NewReader("FROM 1.1.1.1"))
				return script
			},
			exec: func(src *script.Script) error {
				fromCmd, machines, err := exeFrom(src)
				if err != nil {
					return err
				}
				if len(machines) != len(fromCmd.Hosts()) {
					return fmt.Errorf("FROM: expecting %d machines got %d", len(fromCmd.Hosts()), len(machines))
				}
				machine := machines[0]
				if machine.Host() != "1.1.1.1" {
					return fmt.Errorf("FROM machine has unexpected host %s", machine.Host())
				}
				if machine.Port() != "22" {
					return fmt.Errorf("FROM machine has unexpected port %s", machine.Port())
				}

				return nil
			},
		},
		{
			name: "FROM with host:port and global port",
			script: func() *script.Script {
				script, _ := script.Parse(strings.NewReader(`FROM hosts:"1.1.1.1 10.10.10.10:2222" port:2121`))
				return script
			},
			exec: func(src *script.Script) error {
				fromCmd, machines, err := exeFrom(src)
				if err != nil {
					return err
				}
				if len(machines) != len(fromCmd.Hosts()) {
					return fmt.Errorf("FROM: expecting %d machines got %d", len(fromCmd.Hosts()), len(machines))
				}
				m0 := machines[0]
				m1 := machines[1]
				if m0.Host() != "1.1.1.1" || m0.Port() != "2121" {
					return fmt.Errorf("FROM machine0 has unexpected host:port %s:%s", m0.Host(), m0.Port())
				}
				if m1.Host() != "10.10.10.10" || m1.Port() != "2222" {
					return fmt.Errorf("FROM machine1 has unexpected host:port %s:%s", m1.Host(), m1.Port())
				}

				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.exec(test.script()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestExecFROM(t *testing.T) {
	tests := []execTest{
		{
			name: "FROM with multiple addresses",
			source: func() string {
				return `
				ENV host=local
				FROM '$host'
				`
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

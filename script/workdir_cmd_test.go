// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"testing"
)

func TestCommandWORKDIR(t *testing.T) {
	tests := []commandTest{
		{
			name: "WORKDIR",
			command: func(t *testing.T) Directive {
				cmd, err := NewWorkdirCommand(0, "foo/bar")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				wdCmd, ok := c.(*WorkdirCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if wdCmd.Path() != "foo/bar" {
					t.Errorf("WORKDIR has unexpected directory %s", wdCmd.Path())
				}

			},
		},
		{
			name: "WORKDIR/path",
			command: func(t *testing.T) Directive {
				cmd, err := NewWorkdirCommand(0, "path:foo/bar")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				wdCmd, ok := c.(*WorkdirCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if wdCmd.Path() != "foo/bar" {
					t.Errorf("WORKDIR has unexpected directory %s", wdCmd.Path())
				}

			},
		},
		{
			name: "WORKDIR with quoted named param",
			command: func(t *testing.T) Directive {
				cmd, err := NewWorkdirCommand(0, "path:'foo/bar'")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				wdCmd, ok := c.(*WorkdirCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if wdCmd.Path() != "foo/bar" {
					t.Errorf("WORKDIR has unexpected directory %s", wdCmd.Path())
				}

			},
		},
		{
			name: "WORKDIR/expanded vars",
			command: func(t *testing.T) Directive {
				os.Setenv("foopath", "foo/bar")
				cmd, err := NewWorkdirCommand(0, "path:'${foopath}'")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				wdCmd, ok := c.(*WorkdirCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if wdCmd.Path() != "foo/bar" {
					t.Errorf("WORKDIR has unexpected directory %s", wdCmd.Path())
				}

			},
		},
		{
			name: "WORKDIR/multiple args",
			command: func(t *testing.T) Directive {
				cmd, err := NewWorkdirCommand(0, "foo/bar bazz/buzz")
				if err == nil {
					t.Fatal("Expecting error, but got nil")
				}
				return cmd
			},
		},
		{
			name: "WORKDIR/no args",
			command: func(t *testing.T) Directive {
				cmd, err := NewWorkdirCommand(0, "")
				if err == nil {
					t.Fatal("Expecting error, but got nil")
				}
				return cmd
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}

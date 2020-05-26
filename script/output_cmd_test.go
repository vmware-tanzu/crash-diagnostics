// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"testing"
)

func TestCommandOUTPUT(t *testing.T) {
	tests := []commandTest{
		{
			name: "OUTPUT",
			command: func(t *testing.T) Directive {
				cmd, err := NewOutputCommand(0, "foo/bar.tar.gz")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				outCmd, ok := c.(*OutputCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if outCmd.Path() != "foo/bar.tar.gz" {
					t.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}

			},
		},
		{
			name: "OUTPUT/quoted param",
			command: func(t *testing.T) Directive {
				cmd, err := NewOutputCommand(0, "'foo/bar.tar.gz'")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				outCmd, ok := c.(*OutputCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if outCmd.Path() != "foo/bar.tar.gz" {
					t.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}

			},
		},
		{
			name: "OUTPUT/param",
			command: func(t *testing.T) Directive {
				cmd, err := NewOutputCommand(0, "path:foo/bar.tar.gz")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				outCmd, ok := c.(*OutputCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if outCmd.Path() != "foo/bar.tar.gz" {
					t.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}

			},
		},
		{
			name: "OUTPUT/expanded var",
			command: func(t *testing.T) Directive {
				os.Setenv("foopath", "foo/bar.tar.gz")
				cmd, err := NewOutputCommand(0, "$foopath")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				outCmd, ok := c.(*OutputCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if outCmd.Path() != "foo/bar.tar.gz" {
					t.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
				}

			},
		},
		{
			name: "OUTPUT/multiple args",
			command: func(t *testing.T) Directive {
				cmd, err := NewOutputCommand(0, "path:foo/bar path:bazz/buzz")
				if err == nil {
					t.Fatal("Expecting error, but got nil")
				}
				return cmd
			},
		},
		{
			name: "OUTPUT/no args",
			command: func(t *testing.T) Directive {
				cmd, err := NewOutputCommand(0, "OUTPUT")
				if err == nil {
					t.Fatal("Expecting error, but got nil")
				}
				return cmd
			},
		},
		{
			name: "OUTPUT/embedded colon",
			command: func(t *testing.T) Directive {
				cmd, err := NewOutputCommand(0, "path:foo/bar.tar.gz:ignore")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Directive) {
				outCmd, ok := c.(*OutputCommand)
				if !ok {
					t.Errorf("Unexpected type %T in script", c)
				}
				if outCmd.Path() != "foo/bar.tar.gz:ignore" {
					t.Errorf("OUTPUT has unexpected directory %s", outCmd.Path())
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

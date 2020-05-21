// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package script

import (
	"os"
	"testing"
)

func TestCommandCOPY(t *testing.T) {
	tests := []commandTest{
		{
			name: "COPY",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, "/a/b/c")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 1 {
					t.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
				}

				arg := cmd.Paths()[0]
				if arg != "/a/b/c" {
					t.Errorf("COPY has unexpected argument %s", arg)
				}
			},
		},
		{
			name: "COPY/quoted",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, `'/a/b/c'`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 1 {
					t.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
				}

				arg := cmd.Paths()[0]
				if arg != "/a/b/c" {
					t.Errorf("COPY has unexpected argument %s", arg)
				}
			},
		},
		{
			name: "COPY/quoted param",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, `paths:"/a/b/c"`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 1 {
					t.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
				}

				arg := cmd.Paths()[0]
				if arg != "/a/b/c" {
					t.Errorf("COPY has unexpected argument %s", arg)
				}
			},
		},
		{
			name: "COPY/multiple paths",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, "/a/b/c /e/f/g")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 2 {
					t.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
				}
				if cmd.Paths()[0] != "/a/b/c" {
					t.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
				}
				if cmd.Paths()[1] != "/e/f/g" {
					t.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
				}
			},
		},
		{
			name: "COPY/named param",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, "paths:/a/b/c")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 1 {
					t.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
				}

				arg := cmd.Paths()[0]
				if arg != "/a/b/c" {
					t.Errorf("COPY has unexpected argument %s", arg)
				}
			},
		},
		{
			name: "COPY/named param multiple paths",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, `paths:"/a/b/c /e/f/g"`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 2 {
					t.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
				}
				if cmd.Paths()[0] != "/a/b/c" {
					t.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
				}
				if cmd.Paths()[1] != "/e/f/g" {
					t.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
				}
			},
		},
		{
			name: "COPY/var expansion",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, "${foopath1} /e/f/${foodir}")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				os.Setenv("foopath1", "/a/b/c")
				os.Setenv("foodir", "g")
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 2 {
					t.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
				}
				if cmd.Paths()[0] != "/a/b/c" {
					t.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
				}
				if cmd.Paths()[1] != "/e/f/g" {
					t.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
				}
			},
		},
		{
			name: "COPY/no path",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, "")
				if err == nil {
					t.Fatal("Expecting error, but got nil")
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {},
		},
		{
			name: "COPY/colon in path",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, `'/a/:b/c'`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 1 {
					t.Errorf("COPY has unexpected number of paths %d", len(cmd.Paths()))
				}

				arg := cmd.Paths()[0]
				if arg != "/a/:b/c" {
					t.Errorf("COPY has unexpected argument %s", arg)
				}
			},
		},
		{
			name: "COPY/multiple paths with colon",
			command: func(t *testing.T) Command {
				cmd, err := NewCopyCommand(0, `paths:"/a/b/c /e/:f/g"`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, c Command) {
				cmd := c.(*CopyCommand)
				if len(cmd.Paths()) != 2 {
					t.Errorf("COPY has unexpected number of args %d", len(cmd.Paths()))
				}
				if cmd.Paths()[0] != "/a/b/c" {
					t.Errorf("COPY has unexpected argument[0] %s", cmd.Paths()[0])
				}
				if cmd.Paths()[1] != "/e/:f/g" {
					t.Errorf("COPY has unexpected argument[1] %s", cmd.Paths()[1])
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

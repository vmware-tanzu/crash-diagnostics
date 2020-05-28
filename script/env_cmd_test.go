// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"testing"
)

func TestCommandENV(t *testing.T) {
	//tests := []commandTest{
	//	{
	//		name: "ENV",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, "foo=bar")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			envCmd, ok := c.(*EnvCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if len(envCmd.Envs()) != 1 {
	//				t.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
	//			}
	//			env := envCmd.Envs()["foo"]
	//			if env != "bar" {
	//				t.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "ENV/quoted value",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, `foo="bar bazz"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			envCmd := c.(*EnvCommand)
	//			if len(envCmd.Envs()) != 1 {
	//				t.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
	//			}
	//			env := envCmd.Envs()["foo"]
	//			if env != "bar bazz" {
	//				t.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "ENV/named param",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, "vars:abc=defgh")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			envCmd, ok := c.(*EnvCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if len(envCmd.Envs()) != 1 {
	//				t.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
	//			}
	//			env := envCmd.Envs()["abc"]
	//			if env != "defgh" {
	//				t.Errorf("ENV has unexpected value: %#v", envCmd.Envs())
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "ENV/multiple vars",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, `vars:'a=b c=d e=f'`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			envCmd0 := c.(*EnvCommand)
	//			if len(envCmd0.Envs()) != 3 {
	//				t.Errorf("ENV has unexpected number of env %d", len(envCmd0.Envs()))
	//			}
	//			env := envCmd0.Envs()["a"]
	//			if env != "b" {
	//				t.Errorf("ENV has unexpected value a=%s", envCmd0.Envs()["a"])
	//			}
	//			env0, env1 := envCmd0.Envs()["c"], envCmd0.Envs()["e"]
	//			if env0 != "d" || env1 != "f" {
	//				t.Errorf("ENV has unexpected values env[c]=%s and env[e]=%s", env0, env1)
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "ENV/bad format",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, "a=b foo|bar")
	//			if err == nil {
	//				t.Fatal("Expecting failure but got nil")
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {},
	//	},
	//	{
	//		name: "ENV/missing params",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, "")
	//			if err == nil {
	//				t.Fatal("Expecting failure but got nil")
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {},
	//	},
	//	{
	//		name: "ENV/embedded colon",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, "foo=bar:Baz")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			envCmd, ok := c.(*EnvCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if len(envCmd.Envs()) != 1 {
	//				t.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
	//			}
	//			env := envCmd.Envs()["foo"]
	//			if env != "bar:Baz" {
	//				t.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "ENV/quoted embedded colon",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, `foo="bar bazz:bat"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			envCmd := c.(*EnvCommand)
	//			if len(envCmd.Envs()) != 1 {
	//				t.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
	//			}
	//			env := envCmd.Envs()["foo"]
	//			if env != "bar bazz:bat" {
	//				t.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "ENV/multiple embedded colon",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewEnvCommand(0, `vars:'a="b:g" c=d:d e=f'`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			envCmd0 := c.(*EnvCommand)
	//			if len(envCmd0.Envs()) != 3 {
	//				t.Errorf("ENV has unexpected number of env %d", len(envCmd0.Envs()))
	//			}
	//			env := envCmd0.Envs()["a"]
	//			if env != "b:g" {
	//				t.Errorf("ENV has unexpected value a=%s", envCmd0.Envs()["a"])
	//			}
	//			env0, env1 := envCmd0.Envs()["c"], envCmd0.Envs()["e"]
	//			if env0 != "d:d" || env1 != "f" {
	//				t.Errorf("ENV has unexpected values env[c]=%s and env[e]=%s", env0, env1)
	//			}
	//
	//		},
	//	},
	//}
	//
	//for _, test := range tests {
	//	t.Run(test.name, func(t *testing.T) {
	//		runCommandTest(t, test)
	//	})
	//}
}

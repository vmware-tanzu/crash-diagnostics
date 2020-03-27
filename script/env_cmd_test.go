// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"testing"
)

func TestCommandENV(t *testing.T) {
	tests := []commandTest{
		{
			name: "ENV with single arg",
			source: func() string {
				return "ENV foo=bar"
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				if len(envs) != 1 {
					return fmt.Errorf("Script has unexpected number of ENV %d", len(envs))
				}
				envCmd, ok := envs[0].(*EnvCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", envs[0])
				}
				if len(envCmd.Envs()) != 1 {
					return fmt.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
				}
				env := envCmd.Envs()["foo"]
				if env != "bar" {
					return fmt.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
				}
				return nil
			},
		},
		{
			name: "ENV with single quoted arg",
			source: func() string {
				return `ENV foo="bar bazz"`
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				envCmd := envs[0].(*EnvCommand)
				if len(envCmd.Envs()) != 1 {
					return fmt.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
				}
				env := envCmd.Envs()["foo"]
				if env != "bar bazz" {
					return fmt.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
				}
				return nil
			},
		},
		{
			name: "Multiple ENV with multiple args",
			source: func() string {
				return "ENV a=b\nENV 'c=d e=f'"
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				if len(envs) != 2 {
					return fmt.Errorf("Script has unexpected number of ENV %d", len(envs))
				}

				envCmd0, ok := envs[0].(*EnvCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", envs[0])
				}
				if len(envCmd0.Envs()) != 1 {
					return fmt.Errorf("ENV[0] has unexpected number of env %d", len(envCmd0.Envs()))
				}
				env := envCmd0.Envs()["a"]
				if env != "b" {
					return fmt.Errorf("ENV[0] has unexpected value a=%s", envCmd0.Envs()["a"])
				}

				envCmd1, ok := envs[1].(*EnvCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", envs[1])
				}

				if len(envCmd1.Envs()) != 2 {
					return fmt.Errorf("ENV[1] has unexpected number of env %d", len(envCmd1.Envs()))
				}
				env0, env1 := envCmd1.Envs()["c"], envCmd1.Envs()["e"]
				if env0 != "d" || env1 != "f" {
					return fmt.Errorf("ENV[1] has unexpected values env[c]=%s and env[e]=%s", env0, env1)
				}
				return nil
			},
		},
		{
			name: "ENV with named param arg",
			source: func() string {
				return "ENV vars:abc=defgh"
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				if len(envs) != 1 {
					return fmt.Errorf("Script has unexpected number of ENV %d", len(envs))
				}
				envCmd, ok := envs[0].(*EnvCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", envs[0])
				}
				if len(envCmd.Envs()) != 1 {
					return fmt.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
				}
				env := envCmd.Envs()["abc"]
				if env != "defgh" {
					return fmt.Errorf("ENV has unexpected value: %#v", envCmd.Envs())
				}
				return nil
			},
		},
		{
			name: "ENV with multiple quoted vars",
			source: func() string {
				return `ENV vars:'a=b c=d e=f'`
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				if len(envs) != 1 {
					return fmt.Errorf("Script has unexpected number of ENV %d", len(envs))
				}

				envCmd0 := envs[0].(*EnvCommand)
				if len(envCmd0.Envs()) != 3 {
					return fmt.Errorf("ENV has unexpected number of env %d", len(envCmd0.Envs()))
				}
				env := envCmd0.Envs()["a"]
				if env != "b" {
					return fmt.Errorf("ENV has unexpected value a=%s", envCmd0.Envs()["a"])
				}
				env0, env1 := envCmd0.Envs()["c"], envCmd0.Envs()["e"]
				if env0 != "d" || env1 != "f" {
					return fmt.Errorf("ENV has unexpected values env[c]=%s and env[e]=%s", env0, env1)
				}
				return nil
			},
		},
		{
			name: "ENV with var expansion",
			source: func() string {
				return `
				ENV a=b
				ENV 'c=${a} e=${c}'`
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				envCmd0 := envs[0].(*EnvCommand)
				env := envCmd0.Envs()["a"]
				if env != "b" {
					return fmt.Errorf("Unexpected value env[a]=%s", envCmd0.Envs()["a"])
				}

				envCmd1 := envs[1].(*EnvCommand)
				env0, env1 := envCmd1.Envs()["c"], envCmd1.Envs()["e"]
				if env0 != "b" || env1 != "b" {
					return fmt.Errorf("Unexpected values env[c]=%s and env[e]=%s", env0, env1)
				}
				return nil
			},
		},
		{
			name: "ENV with bad formatted values",
			source: func() string {
				return "ENV a=b foo|bar"
			},
			shouldFail: true,
		},
		{
			name: "ENV with missing env",
			source: func() string {
				return "ENV "
			},
			shouldFail: true,
		},
		{
			name: "ENV unquoted with embedded colon",
			source: func() string {
				return "ENV foo=bar:Baz"
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				if len(envs) != 1 {
					return fmt.Errorf("Script has unexpected number of ENV %d", len(envs))
				}
				envCmd, ok := envs[0].(*EnvCommand)
				if !ok {
					return fmt.Errorf("Unexpected type %T in script", envs[0])
				}
				if len(envCmd.Envs()) != 1 {
					return fmt.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
				}
				env := envCmd.Envs()["foo"]
				if env != "bar:Baz" {
					return fmt.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
				}
				return nil
			},
		},
		{
			name: "ENV quoted arg with embedded colon",
			source: func() string {
				return `ENV foo="bar bazz:bat"`
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				envCmd := envs[0].(*EnvCommand)
				if len(envCmd.Envs()) != 1 {
					return fmt.Errorf("ENV has unexpected number of env %d", len(envCmd.Envs()))
				}
				env := envCmd.Envs()["foo"]
				if env != "bar bazz:bat" {
					return fmt.Errorf("ENV has unexpected value: foo=%s", envCmd.Envs()["foo"])
				}
				return nil
			},
		},
		{
			name: "ENV multiple quoted vars with embedded colon",
			source: func() string {
				return `ENV vars:'a="b:g" c=d:d e=f'`
			},
			script: func(s *Script) error {
				envs := s.Preambles[CmdEnv]
				if len(envs) != 1 {
					return fmt.Errorf("Script has unexpected number of ENV %d", len(envs))
				}

				envCmd0 := envs[0].(*EnvCommand)
				if len(envCmd0.Envs()) != 3 {
					return fmt.Errorf("ENV has unexpected number of env %d", len(envCmd0.Envs()))
				}
				env := envCmd0.Envs()["a"]
				if env != "b:g" {
					return fmt.Errorf("ENV has unexpected value a=%s", envCmd0.Envs()["a"])
				}
				env0, env1 := envCmd0.Envs()["c"], envCmd0.Envs()["e"]
				if env0 != "d:d" || env1 != "f" {
					return fmt.Errorf("ENV has unexpected values env[c]=%s and env[e]=%s", env0, env1)
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

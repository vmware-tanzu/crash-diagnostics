// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"testing"
)

func TestCommandKUBECONFIG(t *testing.T) {
	//tests := []commandTest{
	//	{
	//		name: "KUBECONFIG",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeConfigCommand(0, "/a/b/c")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cfg, ok := c.(*KubeConfigCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if cfg.Path() != "/a/b/c" {
	//				t.Errorf("KUBECONFIG has unexpected config %s", cfg.Path())
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "KUBECONFIG/namped param",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeConfigCommand(0, "path:/a/b/c")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cfg, ok := c.(*KubeConfigCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if cfg.Path() != "/a/b/c" {
	//				t.Errorf("KUBECONFIG has unexpected config %s", cfg.Path())
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "KUBECONFIG/quoted named param",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeConfigCommand(0, `path:"/a/b/c"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cfg, ok := c.(*KubeConfigCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if cfg.Path() != "/a/b/c" {
	//				t.Errorf("KUBECONFIG has unexpected config %s", cfg.Path())
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "KUBECONFIG/var expansion",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeConfigCommand(0, `path:$foopath`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			os.Setenv("foopath", "/a/b/c")
	//
	//			cfg, ok := c.(*KubeConfigCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if cfg.Path() != "/a/b/c" {
	//				t.Errorf("KUBECONFIG has unexpected config %s", cfg.Path())
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "KUBECONFIG/embedded colon",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeConfigCommand(0, "/a/:b/c")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cfg, ok := c.(*KubeConfigCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if cfg.Path() != "/a/:b/c" {
	//				t.Errorf("KUBECONFIG has unexpected config %s", cfg.Path())
	//			}
	//
	//		},
	//	},
	//	{
	//		name: "KUBECONFIG/embedded colon param",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeConfigCommand(0, `path:"/a/:b/c"`)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			cfg, ok := c.(*KubeConfigCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if cfg.Path() != "/a/:b/c" {
	//				t.Errorf("KUBECONFIG has unexpected config %s", cfg.Path())
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

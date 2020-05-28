// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"testing"
)

func TestCommandAUTHCONFIG(t *testing.T) {
	//tests := []commandTest{
	//	{
	//		name: "AUTHCONFIG/all params unquoted",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewAuthConfigCommand(0, "username:test-user private-key:/a/b/c")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, cmd Directive) {
	//			authCmd, ok := cmd.(*AuthConfigCommand)
	//			if !ok {
	//				t.Fatalf("Unexpected type %T in script", cmd)
	//			}
	//			if authCmd.GetUsername() != "test-user" {
	//				t.Errorf("Unexpected username %s", authCmd.GetUsername())
	//			}
	//			if authCmd.GetPrivateKey() != "/a/b/c" {
	//				t.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
	//			}
	//		},
	//	},
	//	{
	//		name: "AUTHCONFIG/all params quoted",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewAuthConfigCommand(0, "username:test-user private-key:'/a/b/c'")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, cmd Directive) {
	//			authCmd, ok := cmd.(*AuthConfigCommand)
	//			if !ok {
	//				t.Fatalf("Unexpected type %T in script", cmd)
	//			}
	//			if authCmd.GetUsername() != "test-user" {
	//				t.Errorf("Unexpected username %s", authCmd.GetUsername())
	//			}
	//			if authCmd.GetPrivateKey() != "/a/b/c" {
	//				t.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
	//			}
	//		},
	//	},
	//	{
	//		name: "AUTHCONFIG/only private-key",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewAuthConfigCommand(0, "private-key:/a/b/c")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, cmd Directive) {
	//			authCmd, ok := cmd.(*AuthConfigCommand)
	//			if !ok {
	//				t.Fatalf("Unexpected type %T in script", cmd)
	//			}
	//			if authCmd.GetUsername() != "" {
	//				t.Errorf("Unexpected username %s", authCmd.GetUsername())
	//			}
	//			if authCmd.GetPrivateKey() != "/a/b/c" {
	//				t.Errorf("Unexpected privateKey %s", authCmd.GetPrivateKey())
	//			}
	//		},
	//	},
	//	{
	//		name: "AUTHCONFIG/var expansion",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewAuthConfigCommand(0, "username:${USER} private-key:$fookey")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, cmd Directive) {
	//			os.Setenv("fookey", "/a/b/c")
	//			authCmd := cmd.(*AuthConfigCommand)
	//			if authCmd.GetUsername() != ExpandEnv("$USER") {
	//				t.Errorf("Unexpected username %s", authCmd.GetUsername())
	//			}
	//			if authCmd.GetPrivateKey() != "/a/b/c" {
	//				t.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
	//			}
	//		},
	//	},
	//
	//	{
	//		name: "AUTHCONFIG with bad args",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewAuthConfigCommand(0, "bar private-key:buzz")
	//			if err == nil {
	//				t.Fatalf("Expecting failure but err == nil")
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, cmd Directive) {},
	//	},
	//
	//	{
	//		name: "AUTHCONFIG/embedded colon",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewAuthConfigCommand(0, "username:test-user private-key:'/a/:b/c'")
	//			if err != nil {
	//				t.Error(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, cmd Directive) {
	//			authCmd, ok := cmd.(*AuthConfigCommand)
	//			if !ok {
	//				t.Fatalf("Unexpected type %T in script", cmd)
	//			}
	//			if authCmd.GetUsername() != "test-user" {
	//				t.Errorf("Unexpected username %s", authCmd.GetUsername())
	//			}
	//			if authCmd.GetPrivateKey() != "/a/:b/c" {
	//				t.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
	//			}
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

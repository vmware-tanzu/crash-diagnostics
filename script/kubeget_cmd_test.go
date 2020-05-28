// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"testing"
)

func TestCommandKUBEGET(t *testing.T) {
	//tests := []commandTest{
	//	{
	//		name: "KUBEGET/objects",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeGetCommand(0, "objects")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			kgCmd, ok := c.(*KubeGetCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if kgCmd.What() != "objects" {
	//				t.Errorf("KUBEGET unexpected what: %s", kgCmd.What())
	//			}
	//		},
	//	},
	//	{
	//		name: "KUBEGET/what-param",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeGetCommand(0, "what:logs")
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			kgCmd, ok := c.(*KubeGetCommand)
	//			if !ok {
	//				t.Errorf("Unexpected type %T in script", c)
	//			}
	//			if kgCmd.What() != "logs" {
	//				t.Errorf("KUBEGET unexpected what: %s", kgCmd.What())
	//			}
	//		},
	//	},
	//	{
	//		name: "KUBEGET/all object params",
	//		command: func(t *testing.T) Directive {
	//			cmd, err := NewKubeGetCommand(0,
	//				`objects namespaces:"myns testns" groups:"v1" kinds:"pods events" versions:"1" names:"my-app" labels:"prod" containers:"webapp"`,
	//			)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			return cmd
	//		},
	//		test: func(t *testing.T, c Directive) {
	//			kgCmd := c.(*KubeGetCommand)
	//			if len(kgCmd.Args()) != 8 {
	//				t.Errorf("KUBEGET unexpected param count: %d", len(kgCmd.Args()))
	//			}
	//			// check each param
	//			if kgCmd.What() != "objects" {
	//				t.Errorf("KUBEGET unexpected what: %s", kgCmd.What())
	//			}
	//			if kgCmd.Namespaces() != "myns testns" {
	//				t.Errorf("KUBEGET unexpected namespaces: %s", kgCmd.Namespaces())
	//			}
	//			if kgCmd.Groups() != "v1" {
	//				t.Errorf("KUBEGET unexpected namespaces: %s", kgCmd.Namespaces())
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

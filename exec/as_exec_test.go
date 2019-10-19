// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecAS(t *testing.T) {
	tests := []execTest{
		{
			name: "Exec AS with userid and groupid",
			source: func() string {
				uid := os.Getuid()
				gid := os.Getgid()
				return fmt.Sprintf("AS userid:%d groupid:%d", uid, gid)
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Exec AS with userid only",
			source: func() string {
				uid := os.Getuid()
				return fmt.Sprintf("AS userid:%d", uid)
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Exec AS with expanded vars",
			source: func() string {
				return `AS userid:${USER}`
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Exec AS with unknown uid gid",
			source: func() string {
				return "AS userid:foo"
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

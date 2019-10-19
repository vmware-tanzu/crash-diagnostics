// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func createTestShellScript(t *testing.T, fname string, content string) error {
	execFile, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer execFile.Close()
	t.Logf("Creating shell script file %s", fname)
	_, err = io.Copy(execFile, strings.NewReader(content))
	return err
}
func TestExecENV(t *testing.T) {
	tests := []execTest{
		{
			name: "ENV with with no var expansion",
			source: func() string {
				return "ENV vars:'TEST_A=1 TEST_B=2 TEST_C=3'"
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				if os.Getenv("TEST_A") != "1" {
					t.Errorf("unexpected ENV TEST_A value: %s", os.Getenv("TEST_A"))
				}
				if os.Getenv("TEST_B") != "2" {
					t.Errorf("unexpected ENV TEST_B value: %s", os.Getenv("TEST_B"))
				}
				if os.Getenv("TEST_C") != "3" {
					t.Errorf("unexpected ENV TEST_C value: %s", os.Getenv("TEST_C"))
				}
				return nil
			},
		},
		{
			name: "ENV with chained var expansion",
			source: func() string {
				return `
				ENV vars:'TEST_A=1' 
				ENV vars:'TEST_B=${TEST_A}' 
				ENV 'TEST_C=${USER}'
				`
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				if os.Getenv("TEST_A") != "1" {
					t.Errorf("unexpected ENV TEST_A value: %s", os.Getenv("TEST_A"))
				}
				if os.Getenv("TEST_B") != "1" {
					t.Errorf("unexpected ENV TEST_B value: %s", os.Getenv("TEST_B"))
				}
				if os.Getenv("TEST_C") != os.ExpandEnv("${USER}") {
					t.Errorf("unexpected ENV TEST_C value: %s", os.Getenv("TEST_C"))
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

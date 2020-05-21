// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"testing"

	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

func TestMain(m *testing.M) {
	testcrashd.Init()
	os.Exit(m.Run())
}

type commandTest struct {
	name    string
	command func(*testing.T) Command
	test    func(*testing.T, Command)
}

func runCommandTest(t *testing.T, test commandTest) {
	if test.command == nil {
		t.Fatalf("test %s missing command", test.name)
	}

	if test.test != nil {
		test.test(t, test.command(t))
	}
}

func TestMapArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       string
		expected   ArgMap
		shouldFail bool
	}{
		{
			name:     "MapArgs/single",
			args:     "foo:bar",
			expected: ArgMap{"foo": "bar"},
		},
		{
			name:     "MapArgs/multiple",
			args:     "foo:bar bazz:buzz",
			expected: ArgMap{"foo": "bar", "bazz": "buzz"},
		},
		{
			name:     "MapArgs/quoted",
			args:     `foo:'bar' bazz:"buzz"`,
			expected: ArgMap{"foo": "bar", "bazz": "buzz"},
		},
		{
			name:       "MapArgs/bad format",
			args:       `foo"`,
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			argMap, err := mapArgs(test.args)
			if err != nil && !test.shouldFail {
				t.Fatal(err)
			}
			for k, v := range argMap {
				if test.expected[k] != v {
					t.Fatalf("Unexpected map value: test.expected[%s] = %s, got %s", k, test.expected[k], v)
				}
			}
		})
	}
}

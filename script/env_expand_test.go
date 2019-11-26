// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"os"
	"testing"
)

func TestExpandVarStack(t *testing.T) {
	tests := []struct {
		name  string
		stack func() *runeStack
		test  func(*runeStack)
	}{
		{
			name: "push/pop",
			stack: func() *runeStack {
				s := newRuneStack()
				s.push('a')
				s.push('b')
				s.pop()
				s.push('$')
				return s
			},
			test: func(s *runeStack) {
				if s.depth() != 2 {
					t.Errorf("unexpected stack depth: %d", s.depth())
				}
			},
		},
		{
			name: "push/pop/peek",
			stack: func() *runeStack {
				s := newRuneStack()
				s.push('a')
				s.push('b')
				s.push('$')
				s.push('\\')
				s.pop()
				return s
			},
			test: func(s *runeStack) {
				if s.depth() != 3 {
					t.Errorf("unexpected stack depth: %d", s.depth())
				}
				if s.peek() != '$' {
					t.Errorf("unexpected stack.peek value: %s", string(s.peek()))
				}
			},
		},
		{
			name: "push/pop/isempty",
			stack: func() *runeStack {
				s := newRuneStack()
				s.push('a')
				s.push('b')
				s.pop()
				s.pop()
				s.pop()
				return s
			},
			test: func(s *runeStack) {
				if s.depth() != 0 {
					t.Errorf("unexpected stack.depth: %d", s.depth())
				}
				if !s.isEmpty() {
					t.Errorf("unexpected stack.empty status: %t", s.isEmpty())
				}
				if s.peek() != 0 {
					t.Errorf("unexpected stack.peek value: %s", string(s.peek()))
				}
			},
		},
		{
			name: "push/pop/isempty",
			stack: func() *runeStack {
				s := newRuneStack()
				s.push('a')
				s.push('b')
				s.pop()
				s.pop()
				s.pop()
				s.push('c')
				s.push('d')
				return s
			},
			test: func(s *runeStack) {
				if s.depth() != 2 {
					t.Errorf("unexpected stack.depth: %d", s.depth())
				}
				if s.isEmpty() {
					t.Errorf("unexpected stack.empty status: %t", s.isEmpty())
				}
				if s.peek() != 'd' {
					t.Errorf("unexpected stack.peek value: %s", string(s.peek()))
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.test(test.stack())
		})
	}
}

func TestExpandVar(t *testing.T) {
	tests := []struct {
		name     string
		genStr   func() string
		expected string
	}{
		{
			name:     "no expansion",
			genStr:   func() string { return " Hello, from the world!  " },
			expected: " Hello, from the world!  ",
		},
		{
			name:     `slash - all`,
			genStr:   func() string { return `\\\\\ \\\ \\\` },
			expected: `\\\\\ \\\ \\\`,
		},
		{
			name:     `slash - single middle`,
			genStr:   func() string { return `this \ that` },
			expected: `this \ that`,
		},
		{
			name:     `slash - single end of word`,
			genStr:   func() string { return `this\ that` },
			expected: `this\ that`,
		},
		{
			name:     `slash - single start word`,
			genStr:   func() string { return `this \that` },
			expected: `this \that`,
		},
		{
			name:     `slash - single start of string`,
			genStr:   func() string { return `\this that` },
			expected: `\this that`,
		},
		{
			name:     `slash - single end of string`,
			genStr:   func() string { return `this that\` },
			expected: `this that\`,
		},
		{
			name:     `slash - single inside single word`,
			genStr:   func() string { return `this\that` },
			expected: `this\that`,
		},
		{
			name:     `slash - single inside multi words`,
			genStr:   func() string { return `this w\t t\at` },
			expected: `this w\t t\at`,
		},
		{
			name:     `slash - multi inside single word`,
			genStr:   func() string { return `t\\s that` },
			expected: `t\\s that`,
		},
		{
			name:     `slash - multi inside multi words`,
			genStr:   func() string { return `t\\s t\ha\t` },
			expected: `t\\s t\ha\t`,
		},
		{
			name:     `slash - multi start word`,
			genStr:   func() string { return `this \\\\that` },
			expected: `this \\\\that`,
		},
		{
			name:     `slash - multi middle`,
			genStr:   func() string { return `this \\\\ that` },
			expected: `this \\\\ that`,
		},
		{
			name:     `slash - multi end of word`,
			genStr:   func() string { return `this\\\ that` },
			expected: `this\\\ that`,
		},
		{
			name:     `slash - multi start of string`,
			genStr:   func() string { return `\\\this that` },
			expected: `\\\this that`,
		},
		{
			name:     `slash - multi start of string`,
			genStr:   func() string { return `this that\\\` },
			expected: `this that\\\`,
		},
		{
			name:     `slash - multi inside single word`,
			genStr:   func() string { return `this\\\that` },
			expected: `this\\\that`,
		},
		{
			name:     `escape - start of string`,
			genStr:   func() string { return `\$this that` },
			expected: `$this that`,
		},
		{
			name:     `escape - middle of string`,
			genStr:   func() string { return `this \$is that` },
			expected: `this $is that`,
		},
		{
			name:     `escape - end of string`,
			genStr:   func() string { return `this \$that` },
			expected: `this $that`,
		},
		{
			name:     `escape - with slash at start of string`,
			genStr:   func() string { return `thi\s\ \$that` },
			expected: `thi\s\ $that`,
		},
		{
			name:     `escape - with slash at end of string`,
			genStr:   func() string { return `\$this th\at\` },
			expected: `$this th\at\`,
		},
		{
			name:     `escape - embedded`,
			genStr:   func() string { return `this\$isthat` },
			expected: `this$isthat`,
		},
		{
			name:     `escape - curly vars`,
			genStr:   func() string { return `this \${is} that` },
			expected: `this ${is} that`,
		},
		{
			name:     `escape - curly vars embedded`,
			genStr:   func() string { return `this\${is}that or other` },
			expected: `this${is}that or other`,
		},
		{
			name:     `dollar - all`,
			genStr:   func() string { return `$$$$$ $$ $$$` },
			expected: `$$$$$ $$ $$$`,
		},
		{
			name:     `dollar - single middle`,
			genStr:   func() string { return `foo $ bar` },
			expected: `foo $ bar`,
		},
		{
			name:     `dollar - single end of word`,
			genStr:   func() string { return `foo$ bar` },
			expected: `foo$ bar`,
		},
		{
			name:     `dollar - single end of string`,
			genStr:   func() string { return `foo$ bar$` },
			expected: `foo$ bar$`,
		},
		{
			name:     `var - undeclared var`,
			genStr:   func() string { return `foo $bar` },
			expected: `foo `,
		},
		{
			name: `var - declared at start of string`,
			genStr: func() string {
				os.Setenv("foo", "boo")
				return `$foo bar`
			},
			expected: `boo bar`,
		},
		{
			name: `var - declared at end of string`,
			genStr: func() string {
				os.Setenv("bar", "zaar")
				return `foo $bar`
			},
			expected: `foo zaar`,
		},
		{
			name: `var - embedded`,
			genStr: func() string {
				os.Setenv("bar", "zaar")
				return `foo:$bar:cat`
			},
			expected: `foo:zaar:cat`,
		},
		{
			name: `var - multi embedded`,
			genStr: func() string {
				os.Setenv("bar", "zaar")
				return `foo:$bar:cat:$tar`
			},
			expected: `foo:zaar:cat:`,
		},
		{
			name: `var - multiple declared vars`,
			genStr: func() string {
				os.Setenv("bar", "zaar")
				os.Setenv("bazz", "raaz")
				return `foo $bar with $bazz`
			},
			expected: `foo zaar with raaz`,
		},
		{
			name: `var - multiple declared with missing vars`,
			genStr: func() string {
				os.Setenv("bar", "zaar")
				os.Setenv("bazz", "raaz")
				return `foo ${bar} with $bazz at ${jazz}`
			},
			expected: `foo zaar with raaz at `,
		},
		{
			name: `var - curl vars embedded in words`,
			genStr: func() string {
				os.Setenv("bar", "zaar")
				os.Setenv("bazz", "raaz")
				return `foo${bar}with $bazz at ${jazz}`
			},
			expected: `foozaarwith raaz at `,
		},
		{
			name:     `var - in dollar amount`,
			genStr:   func() string { return `foo $120.00` },
			expected: `foo 20.00`,
		},
		{
			name: `all`,
			genStr: func() string {
				os.Setenv("DIR", "/var/logs")
				return `/bin/bash -c 'files=\$(sudo find $DIR); for f in \$files; do cat \$f; done'`
			},
			expected: `/bin/bash -c 'files=$(sudo find /var/logs); for f in $files; do cat $f; done'`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ExpandEnv(test.genStr())
			if result != test.expected {
				t.Errorf("unexpected expanded result: %s", result)
			}
		})
	}
}

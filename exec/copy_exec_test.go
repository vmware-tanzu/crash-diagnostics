// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecLocalCOPY(t *testing.T) {
	tests := []execTest{
		{
			name: "COPY single files",
			source: func() string {
				return "COPY /tmp/foo0.txt"
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				cpCmd := s.Actions[0].(*script.CopyCommand)
				srcFile := cpCmd.Paths()[0]
				if err := makeTestFakeFile(t, srcFile, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				relPath, err := filepath.Rel("/", srcFile)
				if err != nil {
					return err
				}
				fileName := filepath.Join(workdir.Path(), machine, relPath)
				if _, err := os.Stat(fileName); err != nil {
					return err
				}

				return nil
			},
		},
		{
			name: "COPY multiple files",
			source: func() string {
				return "COPY /tmp/foo0.txt\nCOPY /tmp/foo1.txt /tmp/foo2.txt"
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Paths()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Paths()[0])
				srcFiles = append(srcFiles, cpCmd1.Paths()[1])

				for i, srcFile := range srcFiles {
					if err := makeTestFakeFile(t, srcFile, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
						return err
					}
					defer os.Remove(srcFile)

				}

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				for _, srcFile := range srcFiles {
					relPath, err := filepath.Rel("/", srcFile)
					if err != nil {
						return err
					}
					fileName := filepath.Join(workdir.Path(), machine, relPath)
					if _, err := os.Stat(fileName); err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			name: "COPY directories and files",
			source: func() string {
				return "COPY /tmp/foodir0\nCOPY /tmp/foodir1 /tmp/foo2.txt"
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Paths()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Paths()[0])
				srcFiles = append(srcFiles, cpCmd1.Paths()[1])

				for i, srcFile := range srcFiles {
					if i == 0 || i == 1 {
						if err := makeTestDir(t, srcFile); err != nil {
							return err
						}
						file := filepath.Join(srcFile, fmt.Sprintf("file-%d.txt", i))
						if err := makeTestFakeFile(t, file, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
							return err
						}
					} else {
						if err := makeTestFakeFile(t, srcFile, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
							return err
						}
					}
					defer os.RemoveAll(srcFile)

				}

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				for _, srcFile := range srcFiles {
					relPath, err := filepath.Rel("/", srcFile)
					if err != nil {
						return err
					}
					fileName := filepath.Join(workdir.Path(), machine, relPath)
					if _, err := os.Stat(fileName); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			name: "COPY with var expansion",
			source: func() string {
				os.Setenv("foofile0", "foo0.txt")
				os.Setenv("foofile1", "/tmp/foo1.txt")
				os.Setenv("foo2", "foo2")
				return "COPY /tmp/${foofile0}\nCOPY ${foofile1} /tmp/${foo2}.txt"
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Paths()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Paths()[0])
				srcFiles = append(srcFiles, cpCmd1.Paths()[1])

				for i, srcFile := range srcFiles {
					if err := makeTestFakeFile(t, srcFile, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
						return err
					}
					defer os.Remove(srcFile)

				}

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				for _, srcFile := range srcFiles {
					relPath, err := filepath.Rel("/", srcFile)
					if err != nil {
						return err
					}
					fileName := filepath.Join(workdir.Path(), machine, relPath)
					if _, err := os.Stat(fileName); err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			name: "COPY with globs",
			source: func() string {
				return `
				COPY /tmp/test-dir/*.txt
				COPY /tmp/test-dir/bazz.csv
				`
			},
			exec: func(s *script.Script) error {

				var paths []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				dir := filepath.Dir(cpCmd0.Paths()[0])

				if err := makeTestDir(t, dir); err != nil {
					return fmt.Errorf("failed to crete test dir: %s", err)
				}

				f0 := filepath.Join(dir, "foo.txt")
				if err := makeTestFakeFile(t, f0, "Hello from Foo!"); err != nil {
					return err
				}
				paths = append(paths, f0)

				f1 := filepath.Join(dir, "bar.txt")
				if err := makeTestFakeFile(t, f1, "Hello from Bar!"); err != nil {
					return err
				}
				paths = append(paths, f1)

				f2 := filepath.Join(dir, "bazz.csv")
				if err := makeTestFakeFile(t, f2, "b, a, z, z"); err != nil {
					return err
				}
				paths = append(paths, f2)
				defer os.RemoveAll(dir)

				e := New(s)
				if err := e.Execute(); err != nil {
					return fmt.Errorf("Test command exec failed: %s", err)
				}

				for _, srcFile := range paths {
					if _, err := os.Stat(srcFile); err != nil {
						return fmt.Errorf("Test unable to verify created file: %s", err)
					}
				}
				return nil
			},
		},
		{
			name: "COPY bad source files",
			source: func() string {
				return "COPY /foo/bar.txt"
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

// TestExecRemoteCOPY test COPY command on a remote machine.
// It assumes running account has $HOME/.ssh/id_rsa private key and
// that the remote machine has public key in authorized_keys.
// If setup properly, comment out t.Skip()

func TestExecRemoteCOPY(t *testing.T) {
	t.Skip(`Skipping: test requires an ssh daemon running and a
		passwordless setup using private key specified with AUTHCONFIG command`)

	tests := []execTest{
		{
			name: "COPY single files",
			source: func() string {
				src := `FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				COPY foo.txt`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				cpCmd := s.Actions[0].(*script.CopyCommand)
				srcFile := cpCmd.Paths()[0]
				if err := makeRemoteTestFile(t, machine, srcFile, "HelloFoo"); err != nil {
					return err
				}
				defer removeRemoteTestFile(t, machine, srcFile)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), srcFile)
				if _, err := os.Stat(fileName); err != nil {
					return err
				}

				return nil
			},
		},
		{
			name: "COPY multiple files",
			source: func() string {
				src := `FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				COPY foo0.txt
				COPY foo1.txt foo2.txt`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Paths()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Paths()[0])
				srcFiles = append(srcFiles, cpCmd1.Paths()[1])

				for i, srcFile := range srcFiles {
					if err := makeRemoteTestFile(t, machine, srcFile, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
						return err
					}

					defer removeRemoteTestFile(t, machine, srcFile)
				}

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				for _, srcFile := range srcFiles {
					fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), srcFile)
					if _, err := os.Stat(fileName); err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			name: "COPY directories and files",
			source: func() string {
				src := `FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				COPY foodir0
				COPY foodir1 foo2.txt`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Paths()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Paths()[0])
				srcFiles = append(srcFiles, cpCmd1.Paths()[1])

				for i, srcFile := range srcFiles {
					if i == 0 || i == 1 {
						if err := makeRemoteTestDir(t, machine, srcFile); err != nil {
							return err
						}
						file := filepath.Join(srcFile, fmt.Sprintf("file-%d.txt", i))
						if err := makeRemoteTestFile(t, machine, file, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
							return err
						}
					} else {
						if err := makeRemoteTestFile(t, machine, srcFile, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
							return err
						}
					}
					defer removeRemoteTestFile(t, machine, srcFile)
				}

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				for _, srcFile := range srcFiles {
					fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), srcFile)
					if _, err := os.Stat(fileName); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			name: "COPY with globs",
			source: func() string {
				src := `FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				COPY /tmp/test-dir/*.txt
				COPY /tmp/test-dir/bazz.csv
				`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var paths []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				dir := filepath.Dir(cpCmd0.Paths()[0])

				if err := makeRemoteTestDir(t, machine, dir); err != nil {
					return err
				}

				f0 := filepath.Join(dir, "foo.txt")
				if err := makeRemoteTestFile(t, machine, f0, "Hello from Foo!"); err != nil {
					return err
				}
				paths = append(paths, f0)

				f1 := filepath.Join(dir, "bar.txt")
				if err := makeRemoteTestFile(t, machine, f1, "Hello from Bar!"); err != nil {
					return err
				}
				paths = append(paths, f1)

				f2 := filepath.Join(dir, "bazz.csv")
				if err := makeRemoteTestFile(t, machine, f2, "b, a, z, z"); err != nil {
					return err
				}
				paths = append(paths, f2)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				for _, path := range paths {
					defer removeRemoteTestFile(t, machine, path)
					fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), path)
					if _, err := os.Stat(fileName); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			name: "COPY bad source files",
			source: func() string {
				src := `FROM 127.0.0.1:22
				AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
				COPY foodir0`
				return src
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

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func TestExecCOPY(t *testing.T) {
	tests := []execTest{
		{
			name: "COPY single files",
			source: func() string {
				src := `FROM 127.0.0.1:2222
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
		// {
		// 	name: "COPY multiple files",
		// 	source: func() string {
		// 		src := `FROM 127.0.0.1:2222
		// 		AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
		// 		COPY foo0.txt
		// 		COPY foo1.txt foo2.txt`
		// 		return src
		// 	},
		// 	exec: func(s *script.Script) error {
		// 		machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
		// 		workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

		// 		var srcFiles []string
		// 		cpCmd0 := s.Actions[0].(*script.CopyCommand)
		// 		srcFiles = append(srcFiles, cpCmd0.Paths()[0])
		// 		cpCmd1 := s.Actions[1].(*script.CopyCommand)
		// 		srcFiles = append(srcFiles, cpCmd1.Paths()[0])
		// 		srcFiles = append(srcFiles, cpCmd1.Paths()[1])

		// 		for i, srcFile := range srcFiles {
		// 			if err := makeRemoteTestFile(t, machine, srcFile, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
		// 				return err
		// 			}

		// 			defer removeRemoteTestFile(t, machine, srcFile)
		// 		}

		// 		e := New(s)
		// 		if err := e.Execute(); err != nil {
		// 			return err
		// 		}

		// 		for _, srcFile := range srcFiles {
		// 			fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), srcFile)
		// 			if _, err := os.Stat(fileName); err != nil {
		// 				return err
		// 			}
		// 		}

		// 		return nil
		// 	},
		// },
		// {
		// 	name: "COPY directories and files",
		// 	source: func() string {
		// 		src := `FROM 127.0.0.1:2222
		// 		AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
		// 		COPY foodir0
		// 		COPY foodir1 foo2.txt`
		// 		return src
		// 	},
		// 	exec: func(s *script.Script) error {
		// 		machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
		// 		workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

		// 		var srcFiles []string
		// 		cpCmd0 := s.Actions[0].(*script.CopyCommand)
		// 		srcFiles = append(srcFiles, cpCmd0.Paths()[0])
		// 		cpCmd1 := s.Actions[1].(*script.CopyCommand)
		// 		srcFiles = append(srcFiles, cpCmd1.Paths()[0])
		// 		srcFiles = append(srcFiles, cpCmd1.Paths()[1])

		// 		for i, srcFile := range srcFiles {
		// 			if i == 0 || i == 1 {
		// 				if err := makeRemoteTestDir(t, machine, srcFile); err != nil {
		// 					return err
		// 				}
		// 				file := filepath.Join(srcFile, fmt.Sprintf("file-%d.txt", i))
		// 				if err := makeRemoteTestFile(t, machine, file, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
		// 					return err
		// 				}
		// 			} else {
		// 				if err := makeRemoteTestFile(t, machine, srcFile, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
		// 					return err
		// 				}
		// 			}
		// 			defer removeRemoteTestFile(t, machine, srcFile)
		// 		}

		// 		e := New(s)
		// 		if err := e.Execute(); err != nil {
		// 			return err
		// 		}

		// 		for _, srcFile := range srcFiles {
		// 			fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), srcFile)
		// 			if _, err := os.Stat(fileName); err != nil {
		// 				return err
		// 			}
		// 		}
		// 		return nil
		// 	},
		// },
		// {
		// 	name: "COPY with globs",
		// 	source: func() string {
		// 		src := `FROM 127.0.0.1:2222
		// 		AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
		// 		COPY /tmp/test-dir/*.txt
		// 		COPY /tmp/test-dir/bazz.csv
		// 		`
		// 		return src
		// 	},
		// 	exec: func(s *script.Script) error {
		// 		machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Nodes()[0].Address()
		// 		workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

		// 		var paths []string
		// 		cpCmd0 := s.Actions[0].(*script.CopyCommand)
		// 		dir := filepath.Dir(cpCmd0.Paths()[0])

		// 		if err := makeRemoteTestDir(t, machine, dir); err != nil {
		// 			return err
		// 		}

		// 		f0 := filepath.Join(dir, "foo.txt")
		// 		if err := makeRemoteTestFile(t, machine, f0, "Hello from Foo!"); err != nil {
		// 			return err
		// 		}
		// 		paths = append(paths, f0)

		// 		f1 := filepath.Join(dir, "bar.txt")
		// 		if err := makeRemoteTestFile(t, machine, f1, "Hello from Bar!"); err != nil {
		// 			return err
		// 		}
		// 		paths = append(paths, f1)

		// 		f2 := filepath.Join(dir, "bazz.csv")
		// 		if err := makeRemoteTestFile(t, machine, f2, "b, a, z, z"); err != nil {
		// 			return err
		// 		}
		// 		paths = append(paths, f2)

		// 		e := New(s)
		// 		if err := e.Execute(); err != nil {
		// 			return err
		// 		}

		// 		for _, path := range paths {
		// 			defer removeRemoteTestFile(t, machine, path)
		// 			fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), path)
		// 			if _, err := os.Stat(fileName); err != nil {
		// 				return err
		// 			}
		// 		}
		// 		return nil
		// 	},
		// },
		// {
		// 	name: "COPY bad source files",
		// 	source: func() string {
		// 		src := `FROM 127.0.0.1:2222
		// 		AUTHCONFIG username:${USER} private-key:${HOME}/.ssh/id_rsa
		// 		COPY foodir0`
		// 		return src
		// 	},
		// 	exec: func(s *script.Script) error {
		// 		e := New(s)
		// 		if err := e.Execute(); err != nil {
		// 			return err
		// 		}
		// 		return nil
		// 	},
		// 	shouldFail: true,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}

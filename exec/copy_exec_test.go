package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func TestExecLocalCOPY(t *testing.T) {
	tests := []execTest{
		{
			name: "COPY single files",
			source: func() string {
				return "COPY /tmp/foo0.txt"
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Host()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				cpCmd := s.Actions[0].(*script.CopyCommand)
				srcFile := cpCmd.Args()[0]
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
				fileName := filepath.Join(workdir.Dir(), machine, relPath)
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
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Host()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Args()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Args()[0])
				srcFiles = append(srcFiles, cpCmd1.Args()[1])

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
					fileName := filepath.Join(workdir.Dir(), machine, relPath)
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
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Host()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Args()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Args()[0])
				srcFiles = append(srcFiles, cpCmd1.Args()[1])

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
					fileName := filepath.Join(workdir.Dir(), machine, relPath)
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
		passwordless setup using private key specified by SSHCONFIG command`)

	tests := []execTest{
		{
			name: "COPY single files",
			source: func() string {
				src := `FROM 127.0.0.1:22
				SSHCONFIG {{.Username}}:{{.Home}}/.ssh/id_rsa
				COPY foo.txt`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				cpCmd := s.Actions[0].(*script.CopyCommand)
				srcFile := cpCmd.Args()[0]
				if err := makeRemoteTestFile(t, machine, srcFile, "HelloFoo"); err != nil {
					return err
				}
				defer removeRemoteTestFile(t, machine, srcFile)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fileName := filepath.Join(workdir.Dir(), sanitizeStr(machine), srcFile)
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
				SSHCONFIG {{.Username}}:{{.Home}}/.ssh/id_rsa
				COPY foo0.txt
				COPY foo1.txt foo2.txt`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Args()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Args()[0])
				srcFiles = append(srcFiles, cpCmd1.Args()[1])

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
					fileName := filepath.Join(workdir.Dir(), sanitizeStr(machine), srcFile)
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
				SSHCONFIG {{.Username}}:{{.Home}}/.ssh/id_rsa
				COPY foodir0
				COPY foodir1 foo2.txt`
				return src
			},
			exec: func(s *script.Script) error {
				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Address()
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)

				var srcFiles []string
				cpCmd0 := s.Actions[0].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd0.Args()[0])
				cpCmd1 := s.Actions[1].(*script.CopyCommand)
				srcFiles = append(srcFiles, cpCmd1.Args()[0])
				srcFiles = append(srcFiles, cpCmd1.Args()[1])

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
					fileName := filepath.Join(workdir.Dir(), sanitizeStr(machine), srcFile)
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
				SSHCONFIG {{.Username}}:{{.Home}}/.ssh/id_rsa
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

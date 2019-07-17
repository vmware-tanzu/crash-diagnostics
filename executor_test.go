package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecutor_New(t *testing.T) {
	tests := []struct {
		name   string
		script *script
	}{
		{name: "simple script", script: &script{}},
		{name: "nil script"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := newExecutor(test.script)
			if s.script != test.script {
				t.Error("unexpected script value")
			}
		})
	}
}

func TestExecutor_Exec_Preambles(t *testing.T) {
	tests := []struct {
		name       string
		script     func() *script
		exec       func(*script) error
		shouldFail bool
	}{
		{
			name: "unsupported FROM",
			script: func() *script {
				s, _ := parse(strings.NewReader("FROM foo"))
				return s
			},
			exec: func(s *script) error {
				e := newExecutor(s)
				return e.exec()
			},
			shouldFail: true,
		},
		{
			name: "setup default workdir",
			script: func() *script {
				s, _ := parse(strings.NewReader("FROM local"))
				return s
			},
			exec: func(s *script) error {
				e := newExecutor(s)
				if err := e.exec(); err != nil {
					return err
				}
				if _, err := os.Stat("/tmp/flareout"); err != nil {
					return err
				}
				if err := os.RemoveAll("/tmp/flarout"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "setup workdir /tmp/flarewd",
			script: func() *script {
				s, _ := parse(strings.NewReader("FROM local\nWORKDIR /tmp/flarewd"))
				return s
			},
			exec: func(s *script) error {
				e := newExecutor(s)
				if err := e.exec(); err != nil {
					return err
				}
				if _, err := os.Stat("/tmp/flarewd"); err != nil {
					return err
				}
				if err := os.RemoveAll("/tmp/flarewd"); err != nil {
					return err
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.exec(test.script()); err != nil {
				if !test.shouldFail {
					t.Fatal(err)
				}
				t.Log(err)
				return
			}

		})
	}
}

func TestExecutor_Exec_COPY(t *testing.T) {
	tests := []struct {
		name       string
		script     string
		exec       func(*script) error
		shouldFail bool
	}{
		{
			name:   "copy command with single file",
			script: "FROM local\nCOPY /tmp/flare-foo.txt",
			exec: func(s *script) error {
				workdir := "/tmp/flareout"
				srcFile := s.actions[0].args[0]
				if err := makeTestFakeFile(t, srcFile, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile)

				e := newExecutor(s)
				if err := e.exec(); err != nil {
					return err
				}

				copiedFile := filepath.Join(workdir, filepath.Base(srcFile))
				if _, err := os.Stat(copiedFile); err != nil {
					return err
				}
				if err := os.RemoveAll(workdir); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:   "copy command with multiple files",
			script: "FROM local\nCOPY /tmp/flare-foo.txt /tmp/flare-bar.txt",
			exec: func(s *script) error {
				workdir := "/tmp/flareout"
				srcFile0 := s.actions[0].args[0]
				srcFile1 := s.actions[0].args[1]
				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile0)

				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
					return err
				}
				defer os.Remove(srcFile1)

				e := newExecutor(s)
				if err := e.exec(); err != nil {
					return err
				}

				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
				if _, err := os.Stat(cpFile0); err != nil {
					return err
				}
				if _, err := os.Stat(cpFile1); err != nil {
					return err
				}
				if err := os.RemoveAll(workdir); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:   "copy command with multiple COPYs",
			script: "FROM local\nCOPY /tmp/flare-foo.txt\nCOPY /tmp/flare-bar.txt",
			exec: func(s *script) error {
				workdir := "/tmp/flareout"
				srcFile0 := s.actions[0].args[0]
				srcFile1 := s.actions[1].args[0]
				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile0)

				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
					return err
				}
				defer os.Remove(srcFile1)

				e := newExecutor(s)
				if err := e.exec(); err != nil {
					return err
				}

				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
				if _, err := os.Stat(cpFile0); err != nil {
					return err
				}
				if _, err := os.Stat(cpFile1); err != nil {
					return err
				}
				if err := os.RemoveAll(workdir); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:   "copy command with a directory source",
			script: "FROM local\nCOPY /tmp/flare-src",
			exec: func(s *script) error {
				workdir := "/tmp/flareout"
				srcDir0 := s.actions[0].args[0]
				if err := makeTestDir(t, srcDir0); err != nil {
					return err
				}
				defer os.RemoveAll(srcDir0)
				srcFile0 := filepath.Join(srcDir0, "foo.txt")
				srcFile1 := filepath.Join(srcDir0, "bar.txt")
				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile0)
				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
					return err
				}
				defer os.Remove(srcFile1)

				e := newExecutor(s)
				if err := e.exec(); err != nil {
					return err
				}

				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
				if _, err := os.Stat(cpFile0); err != nil {
					return err
				}
				if _, err := os.Stat(cpFile1); err != nil {
					return err
				}
				if err := os.RemoveAll(workdir); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:   "copy command with a directory source and a file",
			script: "FROM local\nCOPY /tmp/flare-src /tmp/baz.txt",
			exec: func(s *script) error {
				workdir := "/tmp/flareout"
				srcDir0 := s.actions[0].args[0]
				if err := makeTestDir(t, srcDir0); err != nil {
					return err
				}
				defer os.RemoveAll(srcDir0)
				srcFile0 := filepath.Join(srcDir0, "foo.txt")
				srcFile1 := filepath.Join(srcDir0, "bar.txt")
				srcFile2 := s.actions[0].args[1]
				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile0)
				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
					return err
				}
				defer os.Remove(srcFile1)
				if err := makeTestFakeFile(t, srcFile2, "HelloBaz"); err != nil {
					return err
				}
				defer os.Remove(srcFile2)

				e := newExecutor(s)
				if err := e.exec(); err != nil {
					return err
				}

				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
				cpFile2 := filepath.Join(workdir, filepath.Base(srcFile2))
				if _, err := os.Stat(cpFile0); err != nil {
					return err
				}
				if _, err := os.Stat(cpFile1); err != nil {
					return err
				}
				if _, err := os.Stat(cpFile2); err != nil {
					return err
				}
				if err := os.RemoveAll(workdir); err != nil {
					return err
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script, err := parse(strings.NewReader(test.script))
			if err != nil {
				t.Fatal(err)
			}
			if err := test.exec(script); err != nil {
				if !test.shouldFail {
					t.Fatal(err)
				}
				t.Log(err)
				return
			}

		})
	}
}

func makeTestDir(t *testing.T, name string) error {
	if err := os.MkdirAll(name, 0744); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func makeTestFakeFile(t *testing.T, name, content string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	t.Logf("creating test file %s", name)
	_, err = io.Copy(file, strings.NewReader(content))
	return err
}

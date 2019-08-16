package exec

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func TestExecutor_New(t *testing.T) {
	tests := []struct {
		name   string
		script *script.Script
	}{
		{name: "simple script", script: &script.Script{}},
		{name: "nil script"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := New(test.script)
			if s.script != test.script {
				t.Error("unexpected script value")
			}
		})
	}
}

func TestExecutor_Exec_Preambles(t *testing.T) {
	tests := []struct {
		name       string
		script     func() *script.Script
		exec       func(*script.Script) error
		shouldFail bool
	}{
		{
			name: "unsupported FROM",
			script: func() *script.Script {
				s, _ := script.Parse(strings.NewReader("FROM foo"))
				return s
			},
			exec: func(s *script.Script) error {
				defer os.RemoveAll("/tmp/flareout")
				e := New(s)
				return e.Execute()
			},
			shouldFail: true,
		},
		{
			name: "setup default workdir",
			script: func() *script.Script {
				s, _ := script.Parse(strings.NewReader("FROM local"))
				return s
			},
			exec: func(s *script.Script) error {
				defer os.RemoveAll("/tmp/flareout")
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				if _, err := os.Stat("/tmp/flareout"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "setup workdir /tmp/flarewd",
			script: func() *script.Script {
				s, _ := script.Parse(strings.NewReader("FROM local\nWORKDIR /tmp/flarewd"))
				return s
			},
			exec: func(s *script.Script) error {
				defer os.RemoveAll("/tmp/flarewd")
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				if _, err := os.Stat("/tmp/flarewd"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "setup ENV",
			script: func() *script.Script {
				src := "FROM local\nWORKDIR /tmp/flarewd\nENV MSG1=HELLO\nENV MSG2=WORLD MSG3=!\nCAPTURE ./flarec"
				s, _ := script.Parse(strings.NewReader(src))
				return s
			},
			exec: func(s *script.Script) error {
				workdir := s.Preambles[script.CmdWorkDir][0].Args[0]
				defer os.RemoveAll(workdir)
				// create a executable script to apply ENV
				fname := "flarec"
				execFile, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0755)
				if err != nil {
					return err
				}
				t.Logf("Creating test exec file %s", fname)
				_, err = io.Copy(execFile, strings.NewReader("#!/bin/sh\necho $MSG1 $MSG2 $MSG3"))
				if err := execFile.Close(); err != nil {
					return err
				}
				defer os.Remove(fname)

				// test
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				if _, err := os.Stat(workdir); err != nil {
					return err
				}
				fileName := filepath.Join(workdir, fmt.Sprintf("%s.txt", flatCmd("./flarec")))
				file, err := ioutil.ReadFile(fileName)
				if err != nil {
					return err
				}
				if strings.TrimSpace(string(file)) != "HELLO WORLD !" {
					return fmt.Errorf("ENV value not applied during CAPATURE")
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
		exec       func(*script.Script) error
		shouldFail bool
	}{
		{
			name:   "copy command with single file",
			script: "FROM local\nCOPY /tmp/flare-foo.txt",
			exec: func(s *script.Script) error {
				workdir := "/tmp/flareout"
				defer os.RemoveAll(workdir)

				srcFile := s.Actions[0].Args[0]
				if err := makeTestFakeFile(t, srcFile, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				copiedFile := filepath.Join(workdir, filepath.Base(srcFile))
				if _, err := os.Stat(copiedFile); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:   "copy command with multiple files",
			script: "FROM local\nCOPY /tmp/flare-foo.txt /tmp/flare-bar.txt",
			exec: func(s *script.Script) error {
				workdir := "/tmp/flareout"
				defer os.RemoveAll(workdir)

				srcFile0 := s.Actions[0].Args[0]
				srcFile1 := s.Actions[0].Args[1]
				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile0)

				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
					return err
				}
				defer os.Remove(srcFile1)

				e := New(s)
				if err := e.Execute(); err != nil {
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
				return nil
			},
		},
		{
			name:   "copy command with multiple COPYs",
			script: "FROM local\nCOPY /tmp/flare-foo.txt\nCOPY /tmp/flare-bar.txt",
			exec: func(s *script.Script) error {
				workdir := "/tmp/flareout"
				defer os.RemoveAll(workdir)

				srcFile0 := s.Actions[0].Args[0]
				srcFile1 := s.Actions[1].Args[0]
				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
					return err
				}
				defer os.Remove(srcFile0)

				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
					return err
				}
				defer os.Remove(srcFile1)

				e := New(s)
				if err := e.Execute(); err != nil {
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
				return nil
			},
		},
		{
			name:   "copy command with a directory source",
			script: "FROM local\nCOPY /tmp/flare-src",
			exec: func(s *script.Script) error {
				workdir := "/tmp/flareout"
				defer os.RemoveAll(workdir)

				srcDir0 := s.Actions[0].Args[0]
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

				e := New(s)
				if err := e.Execute(); err != nil {
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
				return nil
			},
		},
		{
			name:   "copy command with a directory source and a file",
			script: "FROM local\nCOPY /tmp/flare-src /tmp/baz.txt",
			exec: func(s *script.Script) error {
				workdir := "/tmp/flareout"
				defer os.RemoveAll(workdir)

				srcDir0 := s.Actions[0].Args[0]
				if err := makeTestDir(t, srcDir0); err != nil {
					return err
				}
				defer os.RemoveAll(srcDir0)
				srcFile0 := filepath.Join(srcDir0, "foo.txt")
				srcFile1 := filepath.Join(srcDir0, "bar.txt")
				srcFile2 := s.Actions[0].Args[1]
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

				e := New(s)
				if err := e.Execute(); err != nil {
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
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script, err := script.Parse(strings.NewReader(test.script))
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

func TestExecutor_Exec_CAPTURE(t *testing.T) {
	tests := []struct {
		name       string
		script     func() string
		shouldFail bool
	}{
		{
			name: "capture as default user and usergroup",
			script: func() string {
				return "FROM local\nCAPTURE echo 'helloFoo'"
			},
		},
		{
			name: "capture as specified user only",
			script: func() string {
				uid := os.Getuid()
				return fmt.Sprintf("FROM local\n AS %d\nCAPTURE echo 'hello Bar'", uid)
			},
		},
		{
			name: "capture as specified user and specified group",
			script: func() string {
				uid := os.Getuid()
				gid := os.Getgid()
				return fmt.Sprintf("FROM local\nAS %d:%d\nCAPTURE echo 'hello Bar'", uid, gid)
			},
		},
		{
			name: "capture with badly formatted uid",
			script: func() string {
				return fmt.Sprintf("FROM local\n AS %s\nCAPTURE echo 'hello Bar'", "foo")
			},
			shouldFail: true,
		},
		{
			name: "capture with badly formatted gid",
			script: func() string {
				uid := os.Getuid()
				return fmt.Sprintf("FROM local\n AS %d:%s\nCAPTURE echo 'hello Bar'", uid, "foo")
			},
			shouldFail: true,
		},
		{
			name: "capture with bad permission",
			script: func() string {
				uid := -1
				return fmt.Sprintf("FROM local\n AS %d\nCAPTURE echo 'hello Bar'", uid)
			},
			shouldFail: true,
		},
	}

	exec := func(s *script.Script) error {
		workdir := "/tmp/flareout"
		defer os.RemoveAll(workdir)

		cmdStr := s.Actions[0].Args[0]
		e := New(s)
		if err := e.Execute(); err != nil {
			return err
		}

		fileName := filepath.Join(workdir, fmt.Sprintf("%s.txt", flatCmd(cmdStr)))
		t.Logf("CAPTURE %s -> %s", cmdStr, fileName)
		if _, err := os.Stat(fileName); err != nil {
			return err
		}
		return nil
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script, err := script.Parse(strings.NewReader(test.script()))
			if err != nil {
				t.Fatal(err)
			}
			if err := exec(script); err != nil {
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

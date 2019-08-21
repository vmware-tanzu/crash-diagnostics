package exec

import (
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func TestMain(m *testing.M) {
	loglevel := "debug"
	flag.StringVar(&loglevel, "loglevel", loglevel, "Sets log level")
	flag.Parse()

	if parsed, err := logrus.ParseLevel(loglevel); err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(parsed)
	}
	logrus.SetOutput(os.Stdout)

	os.Exit(m.Run())
}

type execTest struct {
	name       string
	source     func() string
	exec       func(*script.Script) error
	shouldFail bool
}

func runExecutorTest(t *testing.T, test execTest) {
	defer func() {
		if _, err := os.Stat(script.Defaults.WorkdirValue); err != nil {
			t.Log(err)
			return
		}
		if err := os.RemoveAll(script.Defaults.WorkdirValue); err != nil {
			t.Log(err)
		}
	}()

	script, err := script.Parse(strings.NewReader(test.source()))
	if err != nil {
		if !test.shouldFail {
			t.Fatal(err)
		}
		t.Log(err)
		return
	}
	if err := test.exec(script); err != nil {
		if !test.shouldFail {
			t.Fatal(err)
		}
		t.Log(err)
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

// func TestExecutor_Exec_COPY(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		script     string
// 		exec       func(*script.Script) error
// 		shouldFail bool
// 	}{
// 		{
// 			name:   "copy command with single file",
// 			script: "FROM local\nCOPY /tmp/flare-foo.txt",
// 			exec: func(s *script.Script) error {
// 				workdir := "/tmp/flareout"
// 				defer os.RemoveAll(workdir)

// 				srcFile := s.Actions[0].Args[0]
// 				if err := makeTestFakeFile(t, srcFile, "HelloFoo"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile)

// 				e := New(s)
// 				if err := e.Execute(); err != nil {
// 					return err
// 				}

// 				copiedFile := filepath.Join(workdir, filepath.Base(srcFile))
// 				if _, err := os.Stat(copiedFile); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
// 		{
// 			name:   "copy command with multiple files",
// 			script: "FROM local\nCOPY /tmp/flare-foo.txt /tmp/flare-bar.txt",
// 			exec: func(s *script.Script) error {
// 				workdir := "/tmp/flareout"
// 				defer os.RemoveAll(workdir)

// 				srcFile0 := s.Actions[0].Args[0]
// 				srcFile1 := s.Actions[0].Args[1]
// 				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile0)

// 				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile1)

// 				e := New(s)
// 				if err := e.Execute(); err != nil {
// 					return err
// 				}

// 				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
// 				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
// 				if _, err := os.Stat(cpFile0); err != nil {
// 					return err
// 				}
// 				if _, err := os.Stat(cpFile1); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
// 		{
// 			name:   "copy command with multiple COPYs",
// 			script: "FROM local\nCOPY /tmp/flare-foo.txt\nCOPY /tmp/flare-bar.txt",
// 			exec: func(s *script.Script) error {
// 				workdir := "/tmp/flareout"
// 				defer os.RemoveAll(workdir)

// 				srcFile0 := s.Actions[0].Args[0]
// 				srcFile1 := s.Actions[1].Args[0]
// 				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile0)

// 				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile1)

// 				e := New(s)
// 				if err := e.Execute(); err != nil {
// 					return err
// 				}

// 				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
// 				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
// 				if _, err := os.Stat(cpFile0); err != nil {
// 					return err
// 				}
// 				if _, err := os.Stat(cpFile1); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
// 		{
// 			name:   "copy command with a directory source",
// 			script: "FROM local\nCOPY /tmp/flare-src",
// 			exec: func(s *script.Script) error {
// 				workdir := "/tmp/flareout"
// 				defer os.RemoveAll(workdir)

// 				srcDir0 := s.Actions[0].Args[0]
// 				if err := makeTestDir(t, srcDir0); err != nil {
// 					return err
// 				}
// 				defer os.RemoveAll(srcDir0)
// 				srcFile0 := filepath.Join(srcDir0, "foo.txt")
// 				srcFile1 := filepath.Join(srcDir0, "bar.txt")
// 				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile0)
// 				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile1)

// 				e := New(s)
// 				if err := e.Execute(); err != nil {
// 					return err
// 				}

// 				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
// 				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
// 				if _, err := os.Stat(cpFile0); err != nil {
// 					return err
// 				}
// 				if _, err := os.Stat(cpFile1); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
// 		{
// 			name:   "copy command with a directory source and a file",
// 			script: "FROM local\nCOPY /tmp/flare-src /tmp/baz.txt",
// 			exec: func(s *script.Script) error {
// 				workdir := "/tmp/flareout"
// 				defer os.RemoveAll(workdir)

// 				srcDir0 := s.Actions[0].Args[0]
// 				if err := makeTestDir(t, srcDir0); err != nil {
// 					return err
// 				}
// 				defer os.RemoveAll(srcDir0)
// 				srcFile0 := filepath.Join(srcDir0, "foo.txt")
// 				srcFile1 := filepath.Join(srcDir0, "bar.txt")
// 				srcFile2 := s.Actions[0].Args[1]
// 				if err := makeTestFakeFile(t, srcFile0, "HelloFoo"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile0)
// 				if err := makeTestFakeFile(t, srcFile1, "HelloBar"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile1)
// 				if err := makeTestFakeFile(t, srcFile2, "HelloBaz"); err != nil {
// 					return err
// 				}
// 				defer os.Remove(srcFile2)

// 				e := New(s)
// 				if err := e.Execute(); err != nil {
// 					return err
// 				}

// 				cpFile0 := filepath.Join(workdir, filepath.Base(srcFile0))
// 				cpFile1 := filepath.Join(workdir, filepath.Base(srcFile1))
// 				cpFile2 := filepath.Join(workdir, filepath.Base(srcFile2))
// 				if _, err := os.Stat(cpFile0); err != nil {
// 					return err
// 				}
// 				if _, err := os.Stat(cpFile1); err != nil {
// 					return err
// 				}
// 				if _, err := os.Stat(cpFile2); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			script, err := script.Parse(strings.NewReader(test.script))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if err := test.exec(script); err != nil {
// 				if !test.shouldFail {
// 					t.Fatal(err)
// 				}
// 				t.Log(err)
// 				return
// 			}

// 		})
// 	}
// }

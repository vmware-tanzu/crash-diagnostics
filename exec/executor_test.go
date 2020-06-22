// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/parser"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

const (
	testSSHPort = "2222"
)

func TestMain(m *testing.M) {
	testcrashd.Init()
	//
	//sshSvr := testcrashd.NewSSHServer("test-sshd-exec", testSSHPort)
	//logrus.Debug("Attempting to start SSH server")
	//if err := sshSvr.Start(); err != nil {
	//	logrus.Error(err)
	//	os.Exit(1)
	//}
	//
	//testResult := m.Run()
	//
	//logrus.Debug("Stopping SSH server...")
	//if err := sshSvr.Stop(); err != nil {
	//	logrus.Error(err)
	//	os.Exit(1)
	//}
	//
	//os.Exit(testResult)

	// Skipping all tests
	os.Exit(0)
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
		if err := os.RemoveAll(script.Defaults.OutputValue); err != nil {
			t.Log(err)
		}
	}()

	script, err := parser.Parse(strings.NewReader(test.source()))
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
	t.Logf("Making local dir %s", name)
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
	t.Logf("creating local test file %s", name)
	_, err = io.Copy(file, strings.NewReader(content))
	return err
}

func maketTestSSHClient() (*ssh.SSHClient, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	privKey := filepath.Join(usr.HomeDir, ".ssh/id_rsa")
	return ssh.New(usr.Username, privKey, 30), nil
}

func makeRemoteTestFile(t *testing.T, addr, fileName, content string) error {
	sshc, err := maketTestSSHClient()
	if err != nil {
		return err
	}

	if err := sshc.Dial(addr); err != nil {
		return err
	}
	defer sshc.Hangup()

	t.Logf("creating remote test file %s", fileName)
	_, err = sshc.SSHRun(fmt.Sprintf(`echo '%s' > %s`, content, fileName))
	if err != nil {
		return err
	}
	return nil
}

func removeRemoteTestFile(t *testing.T, addr, fileName string) error {
	sshc, err := maketTestSSHClient()
	if err != nil {
		return err
	}

	if err := sshc.Dial(addr); err != nil {
		return err
	}
	defer sshc.Hangup()
	t.Logf("removing remote test file %s", fileName)
	_, err = sshc.SSHRun(fmt.Sprintf("rm -rf %s", fileName))
	if err != nil {
		return err
	}
	return nil
}

func makeRemoteTestDir(t *testing.T, addr, path string) error {
	sshc, err := maketTestSSHClient()
	if err != nil {
		return err
	}

	if err := sshc.Dial(addr); err != nil {
		return err
	}
	defer sshc.Hangup()
	t.Logf("creating remote test  dir %s", path)
	output, err := sshc.SSHRun(fmt.Sprintf("mkdir -p %s", path))
	if err != nil {
		msgBytes, _ := ioutil.ReadAll(output)
		sshErr := fmt.Errorf("ssh command failed: %s: %s", err, string(msgBytes))
		return sshErr
	}
	return nil
}

func getTestFileContent(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
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

func TestExecutor(t *testing.T) {
	tests := []execTest{
		{
			name: "Executing all commands",
			source: func() string {
				var src strings.Builder
				src.WriteString("# This is a sample comment\n")
				src.WriteString("#### START\n")
				src.WriteString(fmt.Sprintf("FROM 127.0.0.1:%s\n", testSSHPort))
				src.WriteString("WORKDIR /tmp/${USER}\n")
				src.WriteString("CAPTURE /bin/echo \"HELLO\"\n")
				src.WriteString("COPY /tmp/buzz.txt\n")
				src.WriteString("ENV MSG0=HELLO MSG1=WORLD BUZZFILE=buzz.txt\n")
				src.WriteString("CAPTURE ./bar.sh\n")
				src.WriteString("COPY /tmp/foodir /tmp/bardir /tmp/${BUZZFILE}\n")
				src.WriteString("##### END")
				return src.String()
			},
			exec: func(s *script.Script) error {
				// create an executable script to apply ENV
				scriptName := "bar.sh"
				sh := `#!/bin/sh
				echo "$MSG1 $MSG2"
				`
				msgExpected := "HELLO WORLD"
				if err := createTestShellScript(t, scriptName, sh); err != nil {
					return fmt.Errorf("failed to create fake shell script bar.sh: %s", err)
				}
				defer os.RemoveAll(scriptName)

				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Hosts()[0]
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				defer os.RemoveAll(workdir.Path())

				// create fake files and dirs to copy
				var srcPaths []string
				for _, cmd := range []script.Command{s.Actions[1], s.Actions[3]} {
					cpCmd := cmd.(*script.CopyCommand)
					for i, path := range cpCmd.Paths() {
						srcPaths = append(srcPaths, path)
						if strings.HasSuffix(path, "dir") { // create dir/file
							if err := makeRemoteTestDir(t, machine, path); err != nil {
								return fmt.Errorf("failed to make test dir %s: %s", path, err)
							}
							file := filepath.Join(path, fmt.Sprintf("file-%d.txt", i))
							if err := makeRemoteTestFile(t, machine, file, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
								return fmt.Errorf("failed to make fake file %s:%s", file, err)
							}
						} else { // create just file
							if err := makeRemoteTestFile(t, machine, path, fmt.Sprintf("HelloFoo-%d", i)); err != nil {
								return fmt.Errorf("failed to make fake file %s: %s", path, err)
							}
						}
						defer os.RemoveAll(path)
					}
				}

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				// validate cap cmds
				for _, cmd := range []script.Command{s.Actions[0], s.Actions[2]} {
					capCmd := cmd.(*script.CaptureCommand)
					fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), fmt.Sprintf("%s.txt", sanitizeStr(capCmd.GetCmdString())))
					if _, err := os.Stat(fileName); err != nil {
						return fmt.Errorf("CAPTURE file validation failed stat for %s: %s", fileName, err)
					}

					if strings.HasSuffix(fileName, ".sh") {
						file, err := ioutil.ReadFile(fileName)
						if err != nil {
							return fmt.Errorf("failed to read fake file %s: %s", file, err)
						}
						if strings.TrimSpace(string(file)) != msgExpected {
							return fmt.Errorf("CAPTURE ./bar.sh generated unexpected content")
						}
					}
				}

				// validate cp cmds
				for _, path := range srcPaths {
					relPath, err := filepath.Rel("/", path)
					if err != nil {
						return err
					}
					fileName := filepath.Join(workdir.Path(), sanitizeStr(machine), relPath)
					if _, err := os.Stat(fileName); err != nil {
						return fmt.Errorf("COPY failed stat file %s: %s", fileName, err)
					}
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

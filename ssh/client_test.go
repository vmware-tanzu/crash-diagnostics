// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

const (
	testSSHPort = "2424"
)

func TestMain(m *testing.M) {
	testcrashd.Init()

	sshSvr := testcrashd.NewSSHServer("test-sshd-sshclient", testSSHPort)
	logrus.Debug("Attempting to start SSH server")
	if err := sshSvr.Start(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	testResult := m.Run()

	logrus.Debug("Stopping SSH server...")
	if err := sshSvr.Stop(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	os.Exit(testResult)
}
func TestSSHClient(t *testing.T) {
	sshHost := fmt.Sprintf("127.0.0.1:%s", testSSHPort)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	privKey := filepath.Join(homeDir, ".ssh/id_rsa")

	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		prepare    func() (*SSHClient, error)
		run        func(*SSHClient) error
		shouldFail bool
	}{
		{
			name: "dial 127.0.0.1",
			prepare: func() (*SSHClient, error) {
				return New(usr.Username, privKey), nil
			},
			run: func(sshClient *SSHClient) error {
				if err := sshClient.Dial(sshHost); err != nil {
					return err
				}
				defer sshClient.Hangup()

				return nil
			},
		},

		{
			name: "ssh run echo hello",
			prepare: func() (*SSHClient, error) {
				return New(usr.Username, privKey), nil
			},
			run: func(sshClient *SSHClient) error {
				if err := sshClient.Dial(sshHost); err != nil {
					return err
				}
				defer sshClient.Hangup()

				reader, err := sshClient.SSHRun("echo 'Hello World!'")
				if err != nil {
					return err
				}
				buff := new(bytes.Buffer)
				if _, err := io.Copy(buff, reader); err != nil {
					return err
				}

				if strings.TrimSpace(buff.String()) != "Hello World!" {
					t.Fatal("SSHRun unexpected result: ", buff.String())
				}
				return nil
			},
		},
		{
			name: "ssh run bad command",
			prepare: func() (*SSHClient, error) {
				return New(usr.Username, privKey), nil
			},
			run: func(sshClient *SSHClient) error {
				if err := sshClient.Dial(sshHost); err != nil {
					return err
				}
				defer sshClient.Hangup()

				if _, err := sshClient.SSHRun("foo bar"); err != nil {
					return err
				}

				return nil
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := test.prepare()

			if err != nil {
				t.Fatal(err)
			}

			if err := test.run(c); err != nil {
				if !test.shouldFail {
					t.Fatal(err)
				}
				t.Log(err)
			}
		})
	}
}

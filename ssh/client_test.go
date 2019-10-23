// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
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

// There is no good/simple way to mock the ssh.Server (without tons of code 😭)
// This test assumes its environment has sshd running configured with a
// known public key (i.e. id_rsa.pub) added to the server's known hosts.
// TODO add script to automate this for testing later.
func TestSSHClient(t *testing.T) {
	sshHost := "127.0.0.1:22"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	privKey := filepath.Join(homeDir, ".ssh/id_rsa")

	t.Skip(fmt.Sprintf(`Skipping: test requires ssh daemon setup at %s for
		passwordless connection using public key for %s`, sshHost, privKey))

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

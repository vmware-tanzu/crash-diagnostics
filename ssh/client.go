// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

// SSHClient represents a client used to connect to an SSH server
type SSHClient struct {
	user       string
	privateKey string
	insecure   bool
	cfg        *ssh.ClientConfig
	sshc       *ssh.Client
	hostKey    ssh.PublicKey
}

// New creates uses the user and privateKeyPath to create an *SSHClient
func New(user string, privateKeyPath string) *SSHClient {
	client := &SSHClient{
		user:       user,
		privateKey: privateKeyPath,
		insecure:   false,
	}
	return client
}

// newInsecure
func newInsecure(user string) *SSHClient {
	client := &SSHClient{
		user:     user,
		insecure: true,
	}
	return client
}

// Dial connects a remote SSH host at address addr
func (c *SSHClient) Dial(addr string) error {
	logrus.Debug("SSH dialing server", addr)

	if c.user == "" {
		return fmt.Errorf("Missing SSH user")
	}

	if !c.insecure {
		logrus.Debug("Connecting using private key file", c.privateKey)
		cfg, err := c.privateKeyConfig()
		if err != nil {
			return err
		}
		c.cfg = cfg
	} else {
		c.cfg = &ssh.ClientConfig{
			User:            c.user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	sshc, err := ssh.Dial("tcp", addr, c.cfg)
	if err != nil {
		return err
	}
	c.sshc = sshc
	return nil
}

// SSHRun executes the specified command on a remote host over SSH
func (c *SSHClient) SSHRun(cmd string, args ...string) (io.Reader, error) {
	cmdStr := strings.TrimSpace(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
	logrus.Debug("Running remote command: ", cmdStr)
	session, err := c.sshc.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	output := new(bytes.Buffer)
	session.Stdout = output
	session.Stderr = output

	if err := session.Run(cmdStr); err != nil {
		return nil, err
	}
	logrus.Debugf("Remote command succeeded: %s", cmdStr)
	return output, nil
}

// Hangup closes the established SSH connection
func (c *SSHClient) Hangup() error {
	return c.sshc.Close()
}

func (c *SSHClient) privateKeyConfig() (*ssh.ClientConfig, error) {
	if _, err := os.Stat(c.privateKey); err != nil {
		return nil, err
	}

	logrus.Debug("Configuring SSH connection with ", c.privateKey)
	key, err := ioutil.ReadFile(c.privateKey)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Found SSH private key ", c.privateKey)

	return &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// not authenticating host
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

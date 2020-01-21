// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/wait"
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

// NewInsecure
func NewInsecure(user string) *SSHClient {
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
		logrus.Debugf("Connecting using private key file %s@%s", c.user, c.privateKey)
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

	// SSH connections with retries
	maxRetries := 30
	retries := wait.Backoff{Steps: maxRetries, Duration: time.Millisecond * 80, Jitter: 0.1}
	if err := wait.ExponentialBackoff(retries, func() (bool, error) {
		sshc, err := ssh.Dial("tcp", addr, c.cfg)
		if err != nil {
			logrus.Errorf("Failed to dial %s (ssh): %s: will retry connection again", addr, err)
			return false, nil
		}
		logrus.Debug("SSH connection establised")
		c.sshc = sshc
		return true, nil
	}); err != nil {
		logrus.Debugf("SSH connection failed after %d tries", maxRetries)
		return err
	}

	return nil
}

// SSHRun executes the specified command on a remote host over SSH
func (c *SSHClient) SSHRun(cmdStr string) (io.Reader, error) {
	logrus.Debug("SSHRun: ", cmdStr)
	session, err := c.sshc.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	output := new(bytes.Buffer)
	session.Stdout = output
	// TODO figure out a way to get output from stderr
	// The following was causing unpredicatiable behavior
	// wher it seems the content would be overwritten by stderr
	// even when stdout returned something.
	//session.Stderr = output

	if err := session.Start(cmdStr); err != nil {
		return nil, err
	}

	if err := session.Wait(); err != nil {
		os.Setenv("CMD_EXITCODE", fmt.Sprintf("%d", 1))
		os.Setenv("CMD_SUCCESS", "false")
		return nil, fmt.Errorf("SSH: error waiting for response: %s", err)
	}

	os.Setenv("CMD_EXITCODE", fmt.Sprintf("%d", 0))
	os.Setenv("CMD_SUCCESS", "true")

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

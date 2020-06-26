// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/echo"
	"k8s.io/apimachinery/pkg/util/wait"
)

type JumpProxyArg struct {
	User string
	Host string
}

type SSHArgs struct {
	User           string
	Host           string
	PrivateKeyPath string
	Port           string
	MaxRetries     int
	JumpProxy      *JumpProxyArg
}

// Run runs a command over SSH and returns the result as a string
func Run(args SSHArgs, cmd string) (string, error) {
	reader, err := sshRunProc(args, cmd)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	if _, err := result.ReadFrom(reader); err != nil {
		return "", err
	}
	return strings.TrimSpace(result.String()), nil
}

// RunRead runs a command over SSH and returns an io.Reader for stdout/stderr
func RunRead(args SSHArgs, cmd string) (io.Reader, error) {
	return sshRunProc(args, cmd)
}

func sshRunProc(args SSHArgs, cmd string) (io.Reader, error) {
	e := echo.New()
	sshCmd, err := makeSSHCmdStr(args)
	if err != nil {
		return nil, err
	}
	effectiveCmd := fmt.Sprintf(`%s "%s"`, sshCmd, cmd)
	logrus.Debug("ssh.run: ", effectiveCmd)

	var proc *echo.Proc
	maxRetries := args.MaxRetries
	if maxRetries == 0 {
		maxRetries = 10
	}
	retries := wait.Backoff{Steps: maxRetries, Duration: time.Millisecond * 80, Jitter: 0.1}
	if err := wait.ExponentialBackoff(retries, func() (bool, error) {
		p := e.RunProc(effectiveCmd)
		if p.Err() != nil {
			logrus.Warn(fmt.Sprintf("ssh: failed to connect to %s: error '%s': retrying connection", args.Host, p.Err()))
			return false, nil
		}
		proc = p
		return true, nil // worked
	}); err != nil {
		logrus.Debugf("ssh.run failed after %d tries", maxRetries)
		return nil, err
	}

	if proc == nil {
		return nil, fmt.Errorf("ssh.run: did get process result")
	}

	return proc.Out(), nil
}

func makeSSHCmdStr(args SSHArgs) (string, error) {
	if args.User == "" {
		return "", fmt.Errorf("SSH: user is required")
	}
	if args.Host == "" {
		return "", fmt.Errorf("SSH: host is required")
	}

	if args.JumpProxy != nil {
		if args.JumpProxy.User == "" || args.JumpProxy.Host == "" {
			return "", fmt.Errorf("SSH: jump user and host are required")
		}
	}

	sshCmdPrefix := func() string {
		if logrus.GetLevel() == logrus.DebugLevel {
			return "ssh -q -o StrictHostKeyChecking=no -v"
		}
		return "ssh -q -o StrictHostKeyChecking=no"
	}

	pkPath := func() string {
		if args.PrivateKeyPath != "" {
			return fmt.Sprintf("-i %s", args.PrivateKeyPath)
		}
		return ""
	}

	port := func() string {
		if args.Port == "" {
			return "-p 22"
		}
		return fmt.Sprintf("-p %s", args.Port)
	}

	jumpProxy := func() string {
		if args.JumpProxy != nil {
			return fmt.Sprintf("-J %s@%s", args.JumpProxy.User, args.JumpProxy.Host)
		}
		return ""
	}
	// build command as
	// ssh -i <pkpath> -P <port> -J <jumpproxy> user@host
	cmd := fmt.Sprintf(
		`%s %s %s %s %s@%s`,
		sshCmdPrefix(), pkPath(), port(), jumpProxy(), args.User, args.Host,
	)
	return cmd, nil
}

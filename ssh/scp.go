// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/echo"
	"k8s.io/apimachinery/pkg/util/wait"
)

// CopyFrom copies one or more files using SCP from remote host
// and returns the paths of files that were successfully copied.
func CopyFrom(args SSHArgs, rootDir string, sourcePath string) error {
	e := echo.New()
	prog := e.Prog.Avail("scp")
	if len(prog) == 0 {
		return fmt.Errorf("scp program not found")
	}

	targetPath := filepath.Join(rootDir, sourcePath)
	targetDir := filepath.Dir(targetPath)
	pathDir, pathFile := filepath.Split(sourcePath)
	if strings.Index(pathFile, "*") != -1 {
		targetPath = filepath.Join(rootDir, pathDir)
		targetDir = targetPath
	}

	if err := os.MkdirAll(targetDir, 0744); err != nil && !os.IsExist(err) {
		return err
	}

	sshCmd, err := makeSCPCmdStr(prog, args, sourcePath)
	if err != nil {
		return fmt.Errorf("scp: failed to build command string: %s", err)
	}

	effectiveCmd := fmt.Sprintf(`%s "%s"`, sshCmd, targetPath)
	logrus.Debug("scp: ", effectiveCmd)

	maxRetries := args.MaxRetries
	if maxRetries == 0 {
		maxRetries = 10
	}
	retries := wait.Backoff{Steps: maxRetries, Duration: time.Millisecond * 80, Jitter: 0.1}
	if err := wait.ExponentialBackoff(retries, func() (bool, error) {
		p := e.RunProc(effectiveCmd)
		if p.Err() != nil {
			logrus.Warn(fmt.Sprintf("scp: failed to connect to %s: error '%s %s': retrying connection", args.Host, p.Err(), p.Result()))
			return false, nil
		}
		return true, nil // worked
	}); err != nil {
		logrus.Debugf("scp failed after %d tries", maxRetries)
		return fmt.Errorf("scp: failed after %d attempt(s): %s", maxRetries, err)
	}

	logrus.Debugf("scp: copied %s", sourcePath)
	return nil
}

func makeSCPCmdStr(progName string, args SSHArgs, sourcePath string) (string, error) {
	if args.User == "" {
		return "", fmt.Errorf("scp: user is required")
	}
	if args.Host == "" {
		return "", fmt.Errorf("scp: host is required")
	}

	if args.ProxyJump != nil {
		if args.ProxyJump.User == "" || args.ProxyJump.Host == "" {
			return "", fmt.Errorf("scp: jump user and host are required")
		}
	}

	scpCmdPrefix := func() string {
		return fmt.Sprintf("%s -rpq -o StrictHostKeyChecking=no", progName)
	}

	pkPath := func() string {
		if args.PrivateKeyPath != "" {
			return fmt.Sprintf("-i %s", args.PrivateKeyPath)
		}
		return ""
	}

	port := func() string {
		if args.Port == "" {
			return "-P 22"
		}
		return fmt.Sprintf("-P %s", args.Port)
	}

	proxyJump := func() string {
		if args.ProxyJump != nil {
			return fmt.Sprintf("-J %s@%s", args.ProxyJump.User, args.ProxyJump.Host)
		}
		return ""
	}
	// build command as
	// scp -i <pkpath> -P <port> -J <proxyjump> user@host:path
	cmd := fmt.Sprintf(
		`%s %s %s %s %s@%s:%s`,
		scpCmdPrefix(), pkPath(), port(), proxyJump(), args.User, args.Host, sourcePath,
	)
	return cmd, nil
}

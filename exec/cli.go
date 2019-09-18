// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	spaceSep = regexp.MustCompile(`\s`)
)

func CliRun(uid, gid uint32, envs []string, cmd string, args ...string) (io.Reader, error) {
	command, output := prepareCmd(cmd, args...)
	command.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true},
	}
	if len(envs) > 0 {
		command.Env = append(os.Environ(), envs...)
	}

	logrus.Debugf("Running %s %v (uid=%d,gid=%d)", cmd, args, uid, gid)
	if err := command.Run(); err != nil {
		return nil, err
	}

	return output, nil
}

func CliParse(cmdStr string) (cmd string, args []string) {
	args = []string{}
	parts := spaceSep.Split(cmdStr, -1)
	if len(parts) == 0 {
		return
	}
	if len(parts) == 1 {
		cmd = parts[0]
		return
	}
	cmd = parts[0]
	args = parts[1:]
	return
}

func prepareCmd(cmd string, args ...string) (*exec.Cmd, io.Reader) {
	output := new(bytes.Buffer)
	command := exec.Command(cmd, args...)
	command.Stdout = output
	command.Stderr = output
	return command, output
}

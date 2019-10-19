// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"bytes"
	"io"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

// CliRun executes specified command using local CLI interface
func CliRun(uid, gid uint32, cmd string, args ...string) (io.Reader, error) {
	command, output := prepareCmd(cmd, args...)
	command.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true},
	}

	logrus.Debugf("Running %s %v (uid=%d,gid=%d)", cmd, args, uid, gid)
	if err := command.Run(); err != nil {
		return nil, err
	}

	return output, nil
}

func prepareCmd(cmd string, args ...string) (*exec.Cmd, io.Reader) {
	output := new(bytes.Buffer)
	command := exec.Command(cmd, args...)
	command.Stdout = output
	command.Stderr = output
	return command, output
}

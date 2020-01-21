// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/echo"
)

// StartSSHServer starts starts sshd process using image linuxserver/openssh-server.DockerRunSSH
// The server opnes up port 2222 with the following command
/*

docker create \
  --name=test-sshd \
  -e PUBLIC_KEY_FILE=$HOME/.ssh/id_rsa.pub \
  -e USER_NAME=$USER \
  -e SUDO_ACCESS=true \
  -p 2222:2222 \
  -v $HOME/.ssh:/config
  linuxserver/openssh-server

*/
func StartSSHServer() error {
	e := echo.New()
	if len(e.Prog.Avail("docker")) == 0 {
		return fmt.Errorf("unable to find docker binary")
	}

	// attemp to stop (if running)
	logrus.Infof("Attempting to stop container %s (if running)", sshContainerName)
	if err := StopSSHServer(); err != nil {
		logrus.Error(err)
	}

	e.SetVar("CONTAINER_NAME", sshContainerName)
	e.SetVar("SSH_PORT", fmt.Sprintf("%s:2222", sshPort))
	e.SetVar("SSH_DOCKER_IMAGE", "vladimirvivien/openssh-server")
	cmd := e.Eval("docker run --rm --detach --name=$CONTAINER_NAME -p $SSH_PORT -e PUBLIC_KEY_FILE=/config/id_rsa.pub -e USER_NAME=$USER -e SUDO_ACCESS=true -v $HOME/.ssh:/config $SSH_DOCKER_IMAGE")
	logrus.Debugf("Starting SSH server: %s", cmd)
	proc := e.RunProc(cmd)
	result := func() string {
		data, _ := ioutil.ReadAll(proc.Out())
		return strings.TrimSpace(string(data))
	}()
	if proc.Err() != nil {
		msg := fmt.Sprintf("%s: %s", proc.Err(), result)
		return fmt.Errorf(msg)
	}
	logrus.Infof("SSH server started: %s", result)

	return nil
}

func StopSSHServer() error {
	e := echo.New()
	e.SetVar("CONTAINER_NAME", sshContainerName)
	proc := e.RunProc("docker stop $CONTAINER_NAME")
	result := func() string {
		data, _ := ioutil.ReadAll(proc.Out())
		return strings.TrimSpace(string(data))
	}()
	if proc.Err() != nil {
		msg := fmt.Sprintf("%s: %s", proc.Err(), result)
		return fmt.Errorf(msg)
	}
	logrus.Info("SSH server stopped: ", result)

	return nil
}

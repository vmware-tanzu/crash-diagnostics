// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

// Executor represents a type that can execute a script
type Executor struct {
	script *script.Script
}

// New returns an *Executor
func New(src *script.Script) *Executor {
	return &Executor{script: src}
}

// Execute executes the configured script
func (e *Executor) Execute() error {
	logrus.Info("Executing script file")

	// exec FROM
	fromCmd, err := exeFrom(e.script)
	if err != nil {
		return err
	}

	// exec WORKDIR
	workdir, err := exeWorkdir(e.script)
	if err != nil {
		return err
	}

	// exec OUTPUT
	output, err := exeOutput(e.script)
	if err != nil {
		return err
	}

	// retrieve KUBECONFIG and setup client connection
	exeClusterInfo(e.script, filepath.Join(workdir.Dir(), "cluster-dump.json"))

	// process actions for each cluster resource specified in FROM
	for _, fromMachine := range fromCmd.Machines() {
		machineWorkdir, err := makeMachineWorkdir(workdir.Dir(), fromMachine)
		if err != nil {
			return err
		}

		switch fromMachine.Host() {
		case "local":
			logrus.Debug("Executing commands on local machine")
			if err := exeLocally(e.script, machineWorkdir); err != nil {
				return err
			}
		default:
			logrus.Debug("Executing remote commands at ", fromMachine.Address())
			if err := exeRemotely(e.script, &fromMachine, machineWorkdir); err != nil {
				return err
			}
		}
	}

	// write result to output
	if err := archiver.Tar(output.Path(), workdir.Dir()); err != nil {
		return err
	}
	logrus.Infof("Created output at path %s", output.Path())
	logrus.Info("Done")

	return nil
}

func makeMachineWorkdir(workdir string, machine script.Machine) (string, error) {
	machineAddr := machine.Address()
	if machineAddr == "local:" {
		machineAddr = machine.Host()
	}
	machineWorkdir := filepath.Join(workdir, sanitizeStr(machineAddr))
	if err := os.MkdirAll(machineWorkdir, 0744); err != nil && !os.IsExist(err) {
		return "", err
	}
	return machineWorkdir, nil
}

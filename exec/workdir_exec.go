// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

// exeWorkdir extract the viable WorkDir command from script, creates
// the working directory if needed, then return the Workdir Command.
func exeWorkdir(src *script.Script) (*script.WorkdirCommand, error) {
	dirs, ok := src.Preambles[script.CmdWorkDir]
	if !ok {
		return nil, fmt.Errorf("Script missing valid %s", script.CmdWorkDir)
	}
	workdir := dirs[0].(*script.WorkdirCommand)
	logrus.Debugf("Using workdir %s", workdir.Path())

	if _, err := os.Stat(workdir.Path()); err != nil {
		if os.IsNotExist(err) {
			logrus.Debugf("Creating  %s", workdir.Path())
			if err := os.MkdirAll(workdir.Path(), 0744); err != nil && !os.IsExist(err) {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return workdir, nil
}

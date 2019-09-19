// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

// exeWorkdir extract the viable WorkDir command from script, creates
// the working directory if needed, then return the Workdir Command.
func exeWorkdir(src *script.Script) (*script.WorkdirCommand, error) {
	dirs, ok := src.Preambles[script.CmdWorkDir]
	if !ok {
		return nil, fmt.Errorf("Script missing valid %s", script.CmdWorkDir)
	}
	workdir := dirs[0].(*script.WorkdirCommand)
	if _, err := os.Stat(workdir.Dir()); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(workdir.Dir(), 0744); err != nil && !os.IsExist(err) {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return workdir, nil
}

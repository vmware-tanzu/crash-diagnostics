// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

// exeOutput extract the OUTPUT command from script, creates
// the output parent directory if needed.
func exeOutput(src *script.Script) (*script.OutputCommand, error) {
	outs, ok := src.Preambles[script.CmdOutput]
	if !ok {
		return nil, fmt.Errorf("Script missing valid %s", script.CmdOutput)
	}

	output := outs[0].(*script.OutputCommand)
	logrus.Debugf("Setting output to %s", output.Path())

	parentPath := filepath.Dir(output.Path())
	if parentPath == "." {
		return output, nil
	}

	// attempt to create parent path
	if _, err := os.Stat(parentPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		logrus.Debugf("Creating directory %s", parentPath)
		if err := os.MkdirAll(parentPath, 0744); err != nil && !os.IsExist(err) {
			return nil, err
		}
	}

	return output, nil
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

// execEnvs saves each declared env variable
// as an ENV for the running process.
func exeEnvs(src *script.Script) error {
	envCmds := src.Preambles[script.CmdEnv]
	for _, envCmd := range envCmds {
		cmd := envCmd.(*script.EnvCommand)
		for name, val := range cmd.Envs() {
			if err := os.Setenv(name, os.ExpandEnv(val)); err != nil {
				return fmt.Errorf("ENV: %s", err)
			}
		}
	}
	return nil
}

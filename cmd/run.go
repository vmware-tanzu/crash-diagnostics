// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/crash-diagnostics/exec"
)

// newRunCommand creates a command to run the Diagnostics script a file
func newRunCommand() *cobra.Command {
	scriptArgs := make(map[string]string)

	cmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "run [flags] <file-name>",
		Short: "Executes a diagnostics script file",
		Long:  "Executes a diagnostics script and collects its output as an archive bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(scriptArgs, args[0])
		},
	}
	cmd.Flags().StringToStringVar(&scriptArgs, "args", scriptArgs, "comma-separated key=value arguments to pass to the diagnostics file")
	return cmd
}

func run(scriptArgs map[string]string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("script file not found: %s", path))
	}

	defer file.Close()

	if err := exec.ExecuteFile(file, scriptArgs); err != nil {
		return errors.Wrap(err, fmt.Sprintf("execution failed for %s", file.Name()))
	}

	return nil
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/crash-diagnostics/exec"
)

type runFlags struct {
	args map[string]string
	file string
}

// newRunCommand creates a command to run the Diaganostics script a file
func newRunCommand() *cobra.Command {
	flags := &runFlags{
		file: "Diagnostics.file",
		args: make(map[string]string),
	}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "run",
		Short: "Executes a diagnostics script file",
		Long:  "Executes a diagnostics script and collects its output as an archive bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(flags)
		},
	}
	cmd.Flags().StringToStringVar(&flags.args, "args", flags.args, "space-separated key=value arguments to passed to diagnostics file")
	cmd.Flags().StringVar(&flags.file, "file", flags.file, "the path to the diagnostics script file to run")
	return cmd
}

func run(flag *runFlags) error {
	file, err := os.Open(flag.file)
	if err != nil {
		return fmt.Errorf("script file not found: %s", flag.file)
	}

	defer file.Close()

	if err := exec.ExecuteFile(file, flag.args); err != nil {
		return fmt.Errorf("execution failed: %s: %s", file.Name(), err)
	}

	return nil
}

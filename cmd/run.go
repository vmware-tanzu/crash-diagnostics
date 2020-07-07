// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/crash-diagnostics/starlark"
)

type runFlags struct {
	file string
	args map[string]string
}

// newRunCommand creates a command to run the Diagnostics script a file
func newRunCommand() *cobra.Command {
	flags := &runFlags{
		file: "Diagnostics.file",
		args: nil,
	}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "run",
		Short: "Executes a diagnostics script file",
		Long:  "Executes a diagnostics script and collects its output as an archive bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(flags, args)
		},
	}
	cmd.Flags().StringVar(&flags.file, "file", flags.file, "the path to the diagnostics script file to run")
	cmd.Flags().StringToStringVar(&flags.args, "args", flags.args, "key-value pair of script arguments which can be used in the diagnostics file")
	return cmd
}

func run(flag *runFlags, _ []string) error {
	if _, err := os.Stat(flag.file); err != nil {
		return fmt.Errorf("unable to find script file %s", flag.file)
	}

	executor := starlark.New()
	if flag.args != nil {
		executor.WithArgs(flag.args)
	}
	if err := executor.Exec(flag.file, nil); err != nil {
		return err
	}
	return nil
}

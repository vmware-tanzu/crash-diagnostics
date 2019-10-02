// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

type runFlags struct {
	file   string
	output string
}

// newRunCommand creates a command to run the Diaganostics script a file
func newRunCommand() *cobra.Command {
	flags := &runFlags{
		file:   "Diagnostics.file",
		output: "out.tar.gz",
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
	cmd.Flags().StringVar(&flags.file, "file", flags.file, "the path to the dianostics script file to run")
	cmd.Flags().StringVar(&flags.output, "output", "", "the path of the generated archive file")
	return cmd
}

func run(flag *runFlags, args []string) error {
	file, err := os.Open(flag.file)
	if err != nil {
		return fmt.Errorf("Unable to find script file %s", flag.file)
	}

	defer file.Close()

	src, err := script.Parse(file)
	if err != nil {
		return err
	}

	// override output if needed
	if flag.output != "" {
		cmd, err := script.NewOutputCommand(0, []string{fmt.Sprintf("path:%s", flag.output)})
		if err != nil {
			return err
		}
		src.Preambles[script.CmdOutput] = []script.Command{cmd}
	}

	exe := exec.New(src)
	if err := exe.Execute(); err != nil {
		return err
	}

	return nil
}

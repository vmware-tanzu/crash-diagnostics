// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/vivienv/crash-diagnostics/exec"
	"gitlab.eng.vmware.com/vivienv/crash-diagnostics/script"
)

type outFlags struct {
	file   string
	output string
}

// out command executes the script and generate a file
// that is compressed into a tarball.
func newOutCommand() *cobra.Command {
	flags := &outFlags{
		file:   "Diagnostics.file",
		output: "out.tar.gz",
	}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "out",
		Short: "outputs an archive from collected data",
		Long:  "outputs an archive from data collected from the specified machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOut(flags, args)
		},
	}
	cmd.Flags().StringVar(&flags.file, "file", flags.file, "the path to the crash-dianostics script file (default ./Dianostics.file)")
	cmd.Flags().StringVar(&flags.output, "output", flags.output, "the path to the generated archive file (default out.tar.gz)")
	return cmd
}

func runOut(flag *outFlags, args []string) error {
	file, err := os.Open(flag.file)
	if err != nil {
		return fmt.Errorf("Unable to find script file %s", flag.file)
	}

	defer file.Close()

	src, err := script.Parse(file)
	if err != nil {
		return err
	}

	// insert output
	if _, ok := src.Preambles[script.CmdOutput]; !ok {
		cmd, err := script.NewOutputCommand(0, []string{fmt.Sprintf("path:%s", flag.output)})
		if err != nil {
			return err
		}
		logrus.Debugf("OUTPUT path:%s (as default)", flag.output)
		src.Preambles[script.CmdOutput] = []script.Command{cmd}
	}

	exe := exec.New(src)
	if err := exe.Execute(); err != nil {
		return err
	}

	return nil
}

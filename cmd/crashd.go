// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/crash-diagnostics/buildinfo"
)

const (
	defaultLogLevel = logrus.InfoLevel
	CliName         = "crashd"
)

// globalFlags flags for the command
type globalFlags struct {
	debug bool
}

// crashDiagnosticsCommand creates a main cli command
func crashDiagnosticsCommand() *cobra.Command {
	flags := &globalFlags{debug: false}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   CliName,
		Short: "runs the crashd program",
		Long:  "Runs the crashd program to execute script that interacts with Kubernetes clusters",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(flags)
		},
		SilenceUsage: true,
		Version:      buildinfo.Version,
	}

	cmd.PersistentFlags().BoolVar(
		&flags.debug,
		"debug",
		flags.debug,
		"sets log level to debug",
	)

	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newBuildinfoCommand())
	return cmd
}

func preRun(flags *globalFlags) error {
	level := defaultLogLevel
	if flags.debug {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)

	return CreateCrashdDir()
}

// Run starts the command
func Run() error {
	logrus.SetOutput(os.Stdout)
	return crashDiagnosticsCommand().Execute()
}

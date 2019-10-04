// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const version = "v0.1.0-alpha.0"
const defaultLogLevel = logrus.InfoLevel

// globalFlags flags for the command
type globalFlags struct {
	debug bool
}

// crashDianosticsCommand creates a main cli command
func crashDiagnosticsCommand() *cobra.Command {
	flags := &globalFlags{debug: false}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "crash-diagnotics",
		Short: "crash-dianostics helps to troubleshoot kubernetes cluster",
		Long:  "crash-diagnotics collects diagnostics from an unresponsive Kubernetes cluster",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(flags)
		},
		SilenceUsage: true,
		Version:      version,
	}

	cmd.PersistentFlags().BoolVar(
		&flags.debug,
		"debug",
		flags.debug,
		"sets log level to debug",
	)

	cmd.AddCommand(newRunCommand())
	return cmd
}

func preRun(flags *globalFlags) error {
	level := defaultLogLevel
	if flags.debug {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)

	return nil
}

// Run satarts the command
func Run() error {
	logrus.SetOutput(os.Stdout)
	return crashDiagnosticsCommand().Execute()
}

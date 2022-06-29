// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/crash-diagnostics/buildinfo"
	"github.com/vmware-tanzu/crash-diagnostics/logging"
)

const (
	defaultLogLevel = logrus.InfoLevel
	CliName         = "crashd"
)

// globalFlags flags for the command
type globalFlags struct {
	debug   bool
	logFile string
}

// crashDiagnosticsCommand creates a main cli command
func crashDiagnosticsCommand() *cobra.Command {
	flags := &globalFlags{debug: false, logFile: "auto"}
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

	cmd.PersistentFlags().StringVar(
		&flags.logFile,
		"log-file",
		flags.logFile,
		"Filepath to log to. Defaults to 'auto' which will generate a unique log file. If empty, will disable logging to a file.",
	)

	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newBuildinfoCommand())
	return cmd
}

func preRun(flags *globalFlags) error {
	if err := CreateCrashdDir(); err != nil {
		return err
	}

	if len(flags.logFile) > 0 {
		// Log everything to file, regardless of settings for CLI.
		filehook, err := logging.NewFileHook(flags.logFile)
		if err != nil {
			logrus.Warning("Failed to log to file, logging to stdout (default)")
		} else {
			logrus.AddHook(filehook)
		}
	}

	level := defaultLogLevel
	if flags.debug {
		level = logrus.DebugLevel
	}
	logrus.AddHook(logging.NewCLIHook(os.Stdout, level))

	// Set to trace so all hooks fire. We will handle levels differently for CLI/file.
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)

	return nil
}

// Run starts the command
func Run() error {
	return crashDiagnosticsCommand().Execute()
}

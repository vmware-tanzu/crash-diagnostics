// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const version = "v0.1.0-alpha.5"
const defaultLogLevel = logrus.WarnLevel

// globalFlags flags for the command
type globalFlags struct {
	logLevel string
}

// crashDianosticsCommand creates a main cli command
func crashDiagnosticsCommand() *cobra.Command {
	flags := &globalFlags{logLevel: defaultLogLevel.String()}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "crash-diagnotics",
		Short: "crash-dianostics helps to investigate an unresponsive kubernetes cluster",
		Long:  "crash-diagnotics collects and analyzes cluster node info from multiple data sources to troubleshoot unresponsive Kubernetes clusters",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(flags)
		},
		SilenceUsage: true,
		Version:      version,
	}

	cmd.PersistentFlags().StringVar(
		&flags.logLevel,
		"loglevel",
		flags.logLevel,
		fmt.Sprintf("log level %v", logrus.AllLevels),
	)

	cmd.AddCommand(newOutCommand())
	return cmd
}

func preRun(flags *globalFlags) error {
	level := defaultLogLevel
	parsed, err := logrus.ParseLevel(flags.logLevel)
	if err != nil {
		logrus.Warnf("Invalid log level [%s], using [%s]", flags.logLevel, level)
	} else {
		level = parsed
	}
	logrus.SetLevel(level)
	return nil
}

// Run satarts the command
func Run() error {
	logrus.SetOutput(os.Stdout)
	return crashDiagnosticsCommand().Execute()
}

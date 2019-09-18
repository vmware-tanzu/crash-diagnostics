// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const version = "v0.1.0-alpha"
const defaultLogLevel = logrus.WarnLevel

// Flags for the kind command
type flareFlags struct {
	logLevel string
}

func flareCommand() *cobra.Command {
	flags := &flareFlags{logLevel: defaultLogLevel.String()}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "flare",
		Short: "flare helps investigate unresponsive kubernetes cluster",
		Long:  "flare helps collect machine info from multiple data sources to troubleshoot Kubernetes 'nodes'",
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

func preRun(flags *flareFlags) error {
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

func Run() error {
	logrus.SetOutput(os.Stdout)
	return flareCommand().Execute()
}

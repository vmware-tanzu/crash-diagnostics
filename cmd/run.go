// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/util"
)

type runFlags struct {
	args     map[string]string
	argsFile string
}

func defaultRunFlags() *runFlags {
	return &runFlags{
		args:     make(map[string]string),
		argsFile: ArgsFile,
	}
}

// newRunCommand creates a command to run the Diagnostics script a file
func newRunCommand() *cobra.Command {
	flags := defaultRunFlags()

	cmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "run <file-name>",
		Short: "Executes a diagnostics script file",
		Long:  "Executes a diagnostics script and collects its output as an archive bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(flags, args[0])
		},
	}
	cmd.Flags().StringToStringVar(&flags.args, "args", flags.args, "comma-separated key=value arguments to pass to the diagnostics file")
	cmd.Flags().StringVar(&flags.argsFile, "args-file", flags.argsFile, "path to the file having key=value arguments to pass to the diagnostics file")
	return cmd
}

func run(flags *runFlags, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "script file not found: %s", path)
	}
	defer file.Close()

	scriptArgs, err := processScriptArguments(flags)
	if err != nil {
		return err
	}

	if err := exec.ExecuteFile(file, scriptArgs); err != nil {
		return errors.Wrapf(err, "execution failed for %s", file.Name())
	}

	return nil
}

// prepares a map of key-value strings to be passed to the execution script
// It builds the map from the args-file as well as the args flag passed to
// the run command.
func processScriptArguments(flags *runFlags) (map[string]string, error) {
	// read inputs from the scriptArgs-file
	scriptArgs, err := util.ReadArgsFile(flags.argsFile)
	if err != nil && flags.argsFile != ArgsFile {
		return nil, errors.Wrapf(err, "failed to parse scriptArgs file: %s", flags.argsFile)
	}

	// values specified by the args flag override the values from the args-file flag
	for k, v := range flags.args {
		scriptArgs[k] = v
	}

	return scriptArgs, nil
}

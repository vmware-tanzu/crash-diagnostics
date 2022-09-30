// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"strings"

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
		Short: "runs a script file",
		Long:  "Parses and executes the specified script file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(flags, args[0])
		},
	}
	cmd.Flags().StringToStringVar(&flags.args, "args", flags.args, "comma-separated key=value pairs passed to the script (i.e. --args 'key0=val0,key1=val1')")
	cmd.Flags().StringVar(&flags.argsFile, "args-file", flags.argsFile, "path to a file containing key=value argument pairs that are passed to the script file")
	return cmd
}

func run(flags *runFlags, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open script file: %s: %w", path, err)
	}
	defer file.Close()

	scriptArgs, err := processScriptArguments(flags)
	if err != nil {
		return err
	}

	if err := exec.ExecuteFile(file, scriptArgs); err != nil {
		return fmt.Errorf("execution failed for %s: %w", file.Name(), err)
	}

	return nil
}

// prepares a map of key-value strings to be passed to the execution script
// It builds the map from the args-file as well as the args flag passed to
// the run command.
func processScriptArguments(flags *runFlags) (map[string]string, error) {
	scriptArgs := map[string]string{}

	// get args from script args file
	err := util.ReadArgsFile(flags.argsFile, scriptArgs)
	if err != nil && flags.argsFile != ArgsFile {
		return nil, fmt.Errorf("failed to parse scriptArgs file %s: %w", flags.argsFile, err)
	}

	// any value specified by the args flag overrides
	// value with same key in the args-file
	for k, v := range flags.args {
		scriptArgs[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}

	return scriptArgs, nil
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/crash-diagnostics/buildinfo"
)

func newBuildinfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "version",
		Short: "prints version",
		Long:  "Prints version information for the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %s\nGitSHA: %s\n", buildinfo.Version, buildinfo.GitSHA)
			return nil
		},
	}
	return cmd
}

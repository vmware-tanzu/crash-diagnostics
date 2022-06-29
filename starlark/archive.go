// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/logging"
	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"
)

// archiveFunc is a built-in starlark function that bundles specified directories into
// an arhive format (i.e. tar.gz)
// Starlark format: archive(output_file=<file name> ,source_paths=list, includeLogs?=[True|False], includeScript?=[True|False])
func archiveFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var outputFile string
	var paths *starlark.List

	// Default to true so that it helps users in debugging.
	includeLogs := true
	includeScript := true

	if err := starlark.UnpackArgs(
		identifiers.archive, args, kwargs,
		"output_file?", &outputFile,
		"source_paths", &paths,
		"includeLogs?", &includeLogs,
		"includeScript?", &includeScript,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.archive, err)
	}

	if len(outputFile) == 0 {
		outputFile = "archive.tar.gz"
	}

	// Always include the script executed and the logs.
	if script := thread.Local(identifiers.scriptName); includeScript && script != nil && len(script.(string)) > 0 {
		if err := paths.Append(starlark.String(script.(string))); err != nil {
			logrus.Warnf("Unexpected error when adding script to archive paths: %v", err)
		}
	}
	if logPath := thread.Local(identifiers.logPath); includeLogs && logPath != nil && len(logPath.(string)) > 0 {
		if err := paths.Append(starlark.String(logPath.(string))); err != nil {
			logrus.Warnf("Unexpected error when adding log path to archive paths: %v", err)
		}
		if err := logging.CloseFileHooks(nil); err != nil {
			logrus.Warnf("Unexpected error when closing file hooks: %v", err)
		}
	}

	if paths != nil && paths.Len() == 0 {
		return starlark.None, fmt.Errorf("%s: one or more paths required", identifiers.archive)
	}

	if err := archiver.Tar(outputFile, getPathElements(paths)...); err != nil {
		return starlark.None, fmt.Errorf("%s failed: %s", identifiers.archive, err)
	}

	return starlark.String(outputFile), nil
}

func getPathElements(paths *starlark.List) []string {
	pathElems := []string{}
	for i := 0; i < paths.Len(); i++ {
		if val, ok := paths.Index(i).(starlark.String); ok {
			pathElems = append(pathElems, string(val))
		}
	}
	return pathElems
}

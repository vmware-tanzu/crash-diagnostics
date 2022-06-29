// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/logging"
	"go.starlark.net/starlark"
)

// logFunc implements a starlark built-in func for simple message logging.
// This iteration uses Go's standard log package.
// Example:
//   log(msg="message", [prefix="info"])
func logFunc(t *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var msg string
	var prefix string
	if err := starlark.UnpackArgs(
		identifiers.log, args, kwargs,
		"msg", &msg,
		"prefix?", &prefix,
	); err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.log, err)
	}

	// retrieve logger from thread
	loggerLocal := t.Local(identifiers.log)
	if loggerLocal == nil {
		addDefaultLogger(t)
		loggerLocal = t.Local(identifiers.log)
	}

	switch logger := loggerLocal.(type) {
	case *log.Logger:
		if prefix != "" {
			logger.Printf("%s: %s", prefix, msg)
		} else {
			logger.Print(msg)
		}
	case *logrus.Logger:
		if prefix != "" {
			logger.Printf("%s: %s", prefix, msg)
		} else {
			logger.Print(msg)
		}
	default:
		return starlark.None, fmt.Errorf("local logger has unknown type %T", loggerLocal)
	}

	return starlark.None, nil
}

func addDefaultLogger(t *starlark.Thread) {
	loggerLocal := t.Local(identifiers.log)
	if loggerLocal == nil {
		logger := logrus.StandardLogger()
		t.SetLocal(identifiers.log, logger)
		if fh := logging.GetFirstFileHook(logger); fh != nil {
			t.SetLocal(identifiers.logPath, fh.FilePath)
		}
	}
}

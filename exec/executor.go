// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/starlark"
)

type ArgMap map[string]string

func Execute(name string, source io.Reader, args ArgMap, restrictedMode bool) error {
	star, err := newExecutor(args, restrictedMode)
	if err != nil {
		return err
	}

	return execute(star, name, source)
}

func ExecuteFile(file *os.File, args ArgMap, restrictedMode bool) error {
	return Execute(file.Name(), file, args, restrictedMode)
}

type StarlarkModule struct {
	Name   string
	Source io.Reader
}

func ExecuteWithModules(name string, source io.Reader, args ArgMap, restrictedMode bool, modules ...StarlarkModule) error {
	star, err := newExecutor(args, restrictedMode)
	if err != nil {
		return err
	}

	// load modules
	for _, mod := range modules {
		if err := star.Preload(mod.Name, mod.Source); err != nil {
			return fmt.Errorf("module load: %w", err)
		}
	}

	return execute(star, name, source)
}

func newExecutor(args ArgMap, restrictedMode bool) (*starlark.Executor, error) {
	star := starlark.New()

	if args != nil {
		starStruct, err := starlark.NewGoValue(args).ToStarlarkStruct("args")
		if err != nil {
			return nil, err
		}

		star.AddPredeclared("args", starStruct)
	}

	return star, nil
}

func execute(star *starlark.Executor, name string, source io.Reader) error {
	if err := star.Exec(name, source); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

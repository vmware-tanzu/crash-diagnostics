// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	starexec "github.com/vmware-tanzu/crash-diagnostics/starlark"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type ArgMap map[string]string

func Execute(name string, source io.Reader, args ArgMap) error {
	star := starexec.New()

	if args != nil {
		starStruct, err := starexec.NewGoValue(args).ToStarlarkStruct("args")
		if err != nil {
			return err
		}

		star.AddPredeclared("args", starStruct)
	}

	err := star.Exec(name, source)
	if err != nil {
		err = fmt.Errorf("exec failed: %w", err)
	}

	return err
}

func ExecuteFile(file *os.File, args ArgMap) error {
	return Execute(file.Name(), file, args)
}

type StarlarkModule struct {
	Name   string
	Source io.Reader
}

func ExecuteWithModules(name string, source io.Reader, args ArgMap, modules ...StarlarkModule) error {
	star := starexec.New()

	if args != nil {
		starStruct, err := starexec.NewGoValue(args).ToStarlarkStruct("args")
		if err != nil {
			return err
		}

		star.AddPredeclared("args", starStruct)
	}

	// load modules
	for _, mod := range modules {
		if err := star.Preload(mod.Name, mod.Source); err != nil {
			return fmt.Errorf("module load: %w", err)
		}
	}

	err := star.Exec(name, source)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

// Run is an alias to Execute which uses the functions package instead.
func Run(name string, source io.Reader, args ArgMap) (starlark.StringDict, error) {
	if args != nil {
		var argsStruct starlarkstruct.Struct
		if err := typekit.Go(args).Starlark(&argsStruct); err != nil {
			return nil, err
		}
		registrar.Register("args", &argsStruct)
	}
	thread := &starlark.Thread{Name: "crashd"}

	if err := setupThreadDefaults(thread); err != nil {
		return nil, fmt.Errorf("thread defaults: %s", err)
	}

	starResult, err := starlark.ExecFile(thread, name, source, registrar.Registry())
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return nil, fmt.Errorf(evalErr.Backtrace())
		}
		return nil, err
	}

	return starResult, err
}

// setupThreadDefaults setups Starlark default thread values
func setupThreadDefaults(thread *starlark.Thread) error {
	if thread == nil {
		return errors.New("thread defaults failed: nil thread")
	}

	if _, err := scriptconf.MakeConfigForThread(thread); err != nil {
		return fmt.Errorf("default script config: failed: %w", err)
	}
	if _, err := sshconf.MakeConfigForThread(thread); err != nil {
		return fmt.Errorf("default ssh config: failed: %w", err)
	}
	return nil
}

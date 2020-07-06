// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"errors"
	"fmt"
	"io"

	"github.com/vladimirvivien/echo"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type Executor struct {
	thread  *starlark.Thread
	predecs starlark.StringDict
	result  starlark.StringDict
}

func New() *Executor {
	return &Executor{
		thread:  &starlark.Thread{Name: "crashd"},
		predecs: newPredeclareds(),
	}
}

func (e *Executor) Exec(name string, source io.Reader) error {
	if err := setupLocalDefaults(e.thread); err != nil {
		return fmt.Errorf("crashd failed: %s", err)
	}

	result, err := starlark.ExecFile(e.thread, name, source, e.predecs)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf(evalErr.Backtrace())
		}
		return err
	}
	e.result = result

	return nil
}

// setupLocalDefaults populates the provided execution thread
// with default configuration values.
func setupLocalDefaults(thread *starlark.Thread) error {
	if thread == nil {
		return errors.New("thread local is nil")
	}
	if err := addDefaultCrashdConf(thread); err != nil {
		return err
	}

	if err := addDefaultSSHConf(thread); err != nil {
		return err
	}

	if err := addDefaultKubeConf(thread); err != nil {
		return err
	}

	return nil
}

// newPredeclareds creates string dictionary containing the
// global built-ins values and functions available to the
// runing script.
func newPredeclareds() starlark.StringDict {
	return starlark.StringDict{
		"os":                          setupOSStruct(),
		identifiers.crashdCfg:         starlark.NewBuiltin(identifiers.crashdCfg, crashdConfigFn),
		identifiers.sshCfg:            starlark.NewBuiltin(identifiers.sshCfg, sshConfigFn),
		identifiers.hostListProvider:  starlark.NewBuiltin(identifiers.hostListProvider, hostListProvider),
		identifiers.resources:         starlark.NewBuiltin(identifiers.resources, resourcesFunc),
		identifiers.run:               starlark.NewBuiltin(identifiers.run, runFunc),
		identifiers.runLocal:          starlark.NewBuiltin(identifiers.runLocal, runLocalFunc),
		identifiers.capture:           starlark.NewBuiltin(identifiers.capture, captureFunc),
		identifiers.captureLocal:      starlark.NewBuiltin(identifiers.capture, captureLocalFunc),
		identifiers.copyFrom:          starlark.NewBuiltin(identifiers.copyFrom, copyFromFunc),
		identifiers.kubeCfg:           starlark.NewBuiltin(identifiers.kubeCfg, kubeConfigFn),
		identifiers.kubeCapture:       starlark.NewBuiltin(identifiers.kubeGet, KubeCaptureFn),
		identifiers.kubeGet:           starlark.NewBuiltin(identifiers.kubeGet, KubeGetFn),
		identifiers.kubeNodesProvider: starlark.NewBuiltin(identifiers.kubeNodesProvider, KubeNodesProviderFn),
	}
}

func kwargsToStringDict(kwargs []starlark.Tuple) (starlark.StringDict, error) {
	if len(kwargs) == 0 {
		return starlark.StringDict{}, nil
	}

	e := echo.New()
	dictionary := make(starlark.StringDict)

	for _, tup := range kwargs {
		key, value := tup[0], tup[1]
		if value.Type() == "string" {
			unquoted := trimQuotes(value.String())
			value = starlark.String(e.Eval(unquoted))
		}
		dictionary[trimQuotes(key.String())] = value
	}

	return dictionary, nil
}

func kwargsToStruct(kwargs []starlark.Tuple) (*starlarkstruct.Struct, error) {
	dict, err := kwargsToStringDict(kwargs)
	if err != nil {
		return &starlarkstruct.Struct{}, err
	}
	return starlarkstruct.FromStringDict(starlarkstruct.Default, dict), nil
}

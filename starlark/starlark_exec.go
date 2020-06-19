// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
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
		thread:  newThreadLocal(),
		predecs: newPredeclareds(),
	}
}

func (e *Executor) Exec(name string, source io.Reader) error {
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

// newThreadLocal creates the execution thread
// and populates default values in the thread.
func newThreadLocal() *starlark.Thread {
	thread := &starlark.Thread{Name: "crashd"}
	addDefaultCrashdConf(thread)
	addDefaultSSHConf(thread)
	addDefaultKubeConf(thread)
	return thread
}

// newPredeclareds creates string dictionary containing the
// global built-ins values and functions available to the
// runing script.
func newPredeclareds() starlark.StringDict {
	return starlark.StringDict{
		"os":                         setupOSStruct(),
		identifiers.crashdCfg:        starlark.NewBuiltin(identifiers.crashdCfg, crashdConfigFn),
		identifiers.sshCfg:           starlark.NewBuiltin(identifiers.sshCfg, sshConfigFn),
		identifiers.hostListProvider: starlark.NewBuiltin(identifiers.hostListProvider, hostListProvider),
		identifiers.resources:        starlark.NewBuiltin(identifiers.resources, resourcesFunc),
		identifiers.run:              starlark.NewBuiltin(identifiers.run, runFunc),
		identifiers.kubeCfg:          starlark.NewBuiltin(identifiers.kubeCfg, kubeConfigFn),
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

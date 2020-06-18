// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"io"

	"github.com/vladimirvivien/echo"
	"go.starlark.net/starlark"
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
	return thread
}

// newPredeclareds creates string dictionary containing the
// global built-ins values and functions available to the
// runing script.
func newPredeclareds() starlark.StringDict {
	return starlark.StringDict{
		identifiers.crashdCfg: starlark.NewBuiltin(identifiers.crashdCfg, crashdConfigFn),
		identifiers.sshCfg:    starlark.NewBuiltin(identifiers.sshCfg, sshConfigFn),
	}
}

func tupleSliceToDict(tuples []starlark.Tuple) (*starlark.Dict, error) {
	if len(tuples) == 0 {
		return &starlark.Dict{}, nil
	}

	dictionary := starlark.NewDict(len(tuples))
	e := echo.New()

	for _, tup := range tuples {
		key, value := tup[0], tup[1]
		if value.Type() == "string" {
			unquoted := trimQuotes(value.String())
			value = starlark.String(e.Eval(unquoted))
		}
		if err := dictionary.SetKey(key, value); err != nil {
			return nil, err
		}
	}

	return dictionary, nil
}

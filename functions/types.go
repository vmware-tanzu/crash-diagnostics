// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package functions

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

// CommandResult represents the result of a command
type CommandResult interface {
	Err() string
	Value() interface{}
}

// DefaultResult is a default result that can be
// used to return result for simple commands
type DefaultResult struct {
	errMsg string
	val    interface{}
}

func NewResult(val interface{}) *DefaultResult {
	return &DefaultResult{val: val}
}
func (c *DefaultResult) AddError(msg string) *DefaultResult {
	c.errMsg = msg
	return c
}

func (c *DefaultResult) Err() string {
	return c.errMsg
}

func (c *DefaultResult) Value() interface{} {
	return c.val
}

func MakeFuncResult(result CommandResult) (starlark.Value, error) {
	dict, err := typekit.GoStructToStringDict(result.Value())
	if err != nil {
		return nil, fmt.Errorf("conversion error: %s", err)
	}
	//dict["Err"] = starlark.String("")
	structVal := starlarkstruct.FromStringDict(starlark.String("CommandResult"), dict)
	fmt.Printf("%T: %#v", structVal, structVal)

	//return starlarkstruct.FromStringDict(starlark.String("CommandResult"), dict), nil
	var star starlarkstruct.Struct
	if err := typekit.Go(result.Value()).Starlark(&star); err != nil {
		return nil, fmt.Errorf("conversion error: %v", err)
	}
	fmt.Printf("%T: %#v", &star, &star)
	return structVal, nil
}

// Command represents a Starlark function that executes
// a specified command and returns a result.
type Command interface {
	Run(*starlark.Thread, interface{}) (CommandResult, error)
}

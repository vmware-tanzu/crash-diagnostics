// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package functions

import (
	"go.starlark.net/starlark"
)

// CommandResult represents the result generated
// by an executable command
type CommandResult interface {
	Err() string
	Value() interface{}
}

// Command represents a Starlark function that executes
// a specified command and returns a result.
type Command interface {
	Run(*starlark.Thread) (CommandResult, error)
}

// Configuration represents Starlark script function that
// collects or generate data to be used as configuration.
type Configuration interface {
	Collect(*starlark.Thread) (interface{}, error)
}

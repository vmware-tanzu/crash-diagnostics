// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package builtins provides the structure for registering
// Starlark functions and other built-in types.
// To avoid cyclic dependencies, use package registrar to manage.
package builtins

import (
	"sync"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"go.starlark.net/starlark"
)

var (
	mutex    sync.Mutex
	registry starlark.StringDict
)

func init() {
	registry = make(starlark.StringDict)
}

// Register registers a Starlark built-in value
func Register(name functions.FunctionName, builtin starlark.Value) {
	mutex.Lock()
	defer mutex.Unlock()
	registry[string(name)] = builtin
}

func Registry() starlark.StringDict {
	return registry
}

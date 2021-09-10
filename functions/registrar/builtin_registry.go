// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package builtins is the façade package that manages the Starlark
// built-in registration of built-ins and other types. The built-in
// Starlark functions are found in package `functions` and the registration
// mechanism is implemented in package functions/builtins.
//
// To register a function, add its import path as a side-effect.
// It's registration will automatically take place.
package registrar

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

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package registrar is the façade package that manages the Starlark
// built-in registration of built-ins and other types. The built-in
// Starlark functions are found in package `functions` and the registration
// mechanism is implemented in package functions/builtins.
//
// To register a function, add its import path as a side-effect.
// It's registration will automatically take place.
package registrar

// Add an import for each function
import (
	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/builtins"
	"go.starlark.net/starlark"
)

// Register registers a Starlark builtin value
func Register(name string, builtin starlark.Value) {
	builtins.Register(functions.FunctionName(name), builtin)
}

func Registry() starlark.StringDict {
	return builtins.Registry()
}

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("resources")
	Func    = resourcesFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// resourcesFunc is a built-in starlark function that provides the convenience
// of returning only the Hosts list of a provider result.  The function takes
// in a value of type provider.Result and returns provider.Result.Resources.
// If an error occurs, this function will halt the script by returning a fatal error.
// For full error-handling control, call the provider directly and inspect its Result.
//
// Args
// - provider: A struct generated from provider.Result
//
// Returns
// - a struct generated from provider.Result.Resources
// - Fatal error if one was returned in provider.Result.Error
//
// Script example: resources(provider=hostlist_provider(hosts=["host1", "host2"]))
func resourcesFunc(thread *starlark.Thread, b *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	// Return fatal error, if any
	if args.ProviderResult.Error != "" {
		return starlark.None, errors.New(args.ProviderResult.Error)
	}

	// Return provider.Result.Resources
	return functions.Result(Name, args.ProviderResult.Resources)
}

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package hostlist

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/registrar"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

var (
	Name    = functions.FunctionName("hostlist_provider")
	Func    = hostListProviderFunc
	Builtin = starlark.NewBuiltin(string(Name), Func)
)

// Register
func init() {
	registrar.Register(Name, Builtin)
}

// hostListProviderFunc is a built-in starlark function that enumerates host resources from a
// provided list of hosts addresses.
//
// Args
// - hosts: list of host addresses (required)
//
// Example: hostlist_provider(hosts=["host1", "host2"])
func hostListProviderFunc(thread *starlark.Thread, b *starlark.Builtin, _ starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var args Args
	if err := typekit.KwargsToGo(kwargs, &args); err != nil {
		return functions.Error(Name, fmt.Errorf("%s: %s", Name, err))
	}

	result := Run(thread, args)

	// convert and return result
	return functions.Result(Name, result)
}

// Run runs the command function
func Run(t *starlark.Thread, args Args) providers.Result {
	if len(args.Hosts) == 0 {
		return providers.Result{Error: "host list is required"}
	}

	return providers.Result{
		Resources: providers.Resources{
			Provider: string(Name),
			Hosts:    args.Hosts,
		},
	}
}

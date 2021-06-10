// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package copy_from represents the `copy_from` Starlark function.
package copy_from

type Args struct {
	Path      string   `name:"path"`
	Resources []string `name:"resources" optional:"true"`
	Workdir   string   `name:"workdir" optional:"true"`
}

type Result struct {
	Error string `name:"error"`
}

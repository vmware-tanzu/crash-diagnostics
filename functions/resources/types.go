// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package resources represents the `resources` starlark function
package resources

type Args struct {
	Hosts     []string   `name:"hosts" optional:"true"`
	Provider  []string `name:"resources" optional:"true"`
}

type Result struct {
	Error string `name:"error"`
}

// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package hostlist_provider represents a resource provider created by  '
// starlark function `hostlist_provider()`
package hostlist_provider

type Args struct {
	Hosts []string `name:"hosts"`
}
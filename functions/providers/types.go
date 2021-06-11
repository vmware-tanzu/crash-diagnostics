// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package providers contains implementation of resource providers
// that are used to enumerate compute resources to be used in
// script functions.
package providers

var (
	ResourcesIdentifier = "resources"
)

type Resources struct {
	Error    string   `name:"error"`
	Provider string   `name:"provider"`
	Hosts    []string `name:"hosts"`
}

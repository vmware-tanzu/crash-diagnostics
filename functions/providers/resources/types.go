// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package resources implements the `resources` script function which
// is a convenience wrapper that takes a Resources type ane returns
// the Hosts field as a value.
package resources

import (
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
)

type Args struct {
	ProviderResult providers.Result `name:"provider"`
}

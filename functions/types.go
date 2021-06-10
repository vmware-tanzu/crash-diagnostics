// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package functions

type FunctionName string

type ProviderResources struct {
	Error string   `name:"error"`
	Hosts []string `name:"hosts"`
}


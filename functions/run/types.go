// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package run

// LocalProc represents the result of executing a local process
// from a Starlark script.
type LocalProc struct {
	Pid      int64
	Result   string
	ExitCode int64
}

// LocalProc represents the result of executing a remote process
// from a Starlark script.
type RemoteProc struct {
	resource string
	Result   string
	ExitCode int64
}

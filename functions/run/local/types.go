// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package local

type Args struct {
	Cmd string `name:"cmd"`
}

type Result struct {
	Error    string `name:"error"`
	Pid      int64  `name:"pid"`
	Result   string `name:"result"`
	ExitCode int64  `name:"exit_code"`
}

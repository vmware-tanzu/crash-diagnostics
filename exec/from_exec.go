// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func exeFrom(src *script.Script) (*script.FromCommand, error) {
	fromCmds, ok := src.Preambles[script.CmdFrom]
	if !ok {
		return nil, fmt.Errorf("Script missing valid %s", script.CmdFrom)
	}
	if len(fromCmds) < 1 {
		return nil, fmt.Errorf("Script missing valid %s", script.CmdFrom)
	}

	return fromCmds[0].(*script.FromCommand), nil
}

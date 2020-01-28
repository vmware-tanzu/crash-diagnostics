// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"strings"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

func exeFrom(src *script.Script) (*script.FromCommand, []*script.Machine, error) {
	fromCmds, ok := src.Preambles[script.CmdFrom]
	if !ok {
		return nil, nil, fmt.Errorf("%s not defined", script.CmdFrom)
	}
	if len(fromCmds) < 1 {
		return nil, nil, fmt.Errorf("script missing valid %s", script.CmdFrom)
	}

	fromCmd, ok := fromCmds[0].(*script.FromCommand)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected type %T for %s", fromCmd, script.CmdFrom)
	}

	var machines []*script.Machine
	// retrieve from host list
	for _, host := range fromCmd.Hosts() {
		var addr, port, name string
		parts := strings.Split(host, ":")
		if len(parts) > 1 {
			addr = parts[0]
			port = parts[1]
			name = host
		} else {
			addr = parts[0]
			port = fromCmd.Port()
			name = host
		}
		machines = append(machines, script.NewMachine(addr, port, name))
	}

	return fromCmd, machines, nil
}

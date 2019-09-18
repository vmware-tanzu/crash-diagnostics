// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"net"
)

// Machine represents a source machine
type Machine struct {
	host string
	port string
}

func NewMachine(host, port string) *Machine {
	return &Machine{host: host, port: port}
}

func (m *Machine) Address() string {
	return net.JoinHostPort(m.host, m.port)
}

func (m *Machine) Host() string {
	return m.host
}

func (m *Machine) Port() string {
	return m.port
}

// FromCommand represents instruction
// FROM <source>
// Where source can be:
// 1. Local : the current machine
// 2. List of machine names/addresses
// 3. cluster: uses Kuberentes cluster information to get list
type FromCommand struct {
	cmd
	machines []Machine
}

// NewFromCommand creates a value of type FromCommand
func NewFromCommand(index int, args []string) (*FromCommand, error) {
	cmd := &FromCommand{cmd: cmd{index: index, name: CmdFrom, args: args}}

	if err := validateCmdArgs(CmdFrom, args); err != nil {
		return nil, err
	}

	for _, arg := range args {
		var host, port string
		switch {
		case arg == "local":
			host = "local"
		case arg == "cluster":
			host = "cluster"
		default:
			h, p, err := net.SplitHostPort(arg)
			if err != nil {
				return nil, fmt.Errorf("FROM command: %s", err)
			}
			host = h
			port = p
			if p == "" {
				port = "22"
			}
		}
		cmd.machines = append(cmd.machines, *NewMachine(host, port))
	}

	return cmd, nil
}

func (c *FromCommand) Index() int {
	return c.cmd.index
}

func (c *FromCommand) Name() string {
	return c.cmd.name
}

func (c *FromCommand) Args() []string {
	return c.cmd.args
}

func (c *FromCommand) Machines() []Machine {
	return c.machines
}

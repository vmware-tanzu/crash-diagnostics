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

// NewMachine returns a new machine
func NewMachine(host, port string) *Machine {
	return &Machine{host: host, port: port}
}

// Address returns the host:port address
func (m *Machine) Address() string {
	return net.JoinHostPort(m.host, m.port)
}

// Host returns the host of the address
func (m *Machine) Host() string {
	return m.host
}

// Port returns the port of the address
func (m *Machine) Port() string {
	return m.port
}

// FromCommand represents FROM directive:
//
// FROM host0:port host1:port...hostN:port
//
// Each host must be specified as an address with host:port.
type FromCommand struct {
	cmd
	machines []Machine
}

// NewFromCommand parses the args and returns *FromCommand
func NewFromCommand(index int, rawArgs string) (*FromCommand, error) {
	if err := validateRawArgs(CmdFrom, rawArgs); err != nil {
		return nil, err
	}

	// by default hosts will store host addresses
	argMap := map[string]string{"hosts": rawArgs}
	cmd := &FromCommand{cmd: cmd{index: index, name: CmdFrom, args: argMap}}
	if err := validateCmdArgs(CmdFrom, argMap); err != nil {
		return nil, err
	}

	// populate machine representations
	for _, arg := range spaceSep.Split(rawArgs, -1) {
		var host, port string
		switch {
		case arg == "local":
			host = "local"
		case arg == "cluster":
			host = "cluster"
		default:
			h, p, err := net.SplitHostPort(arg)
			if err != nil {
				return nil, fmt.Errorf("FROM: %s", err)
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

// Index is the position of the command in the script
func (c *FromCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *FromCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *FromCommand) Args() map[string]string {
	return c.cmd.args
}

// Machines returns a slice of Machines to which to connect
func (c *FromCommand) Machines() []Machine {
	return c.machines
}

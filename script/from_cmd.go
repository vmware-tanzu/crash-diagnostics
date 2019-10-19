// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Machine represents a machine as defined in FROM
type Machine struct {
	address string
}

// NewMachine returns a new *Machine
func NewMachine(addr string) *Machine {
	return &Machine{address: addr}
}

// Address returns the host:port address
func (m *Machine) Address() string {
	return m.address
}

// Host returns the host of the address
func (m *Machine) Host() (string, error) {
	h, _, err := net.SplitHostPort(m.Address())
	if err != nil {
		return "", err
	}
	return h, nil
}

// Port returns the port of the address
func (m *Machine) Port() (string, error) {
	_, p, err := net.SplitHostPort(m.Address())
	if err != nil {
		return "", err
	}
	return p, nil
}

// FromCommand represents FROM directive which may take
// one of the following forms:
//
//     FROM host0:port host1:port ... hostN:port
//     FROM hosts:"host0:port host1:port ... hostN:port"
//
// Each host must be specified as an address endpoint with host:port.
type FromCommand struct {
	cmd
	machines []Machine
}

// NewFromCommand parses the args and returns *FromCommand
func NewFromCommand(index int, rawArgs string) (*FromCommand, error) {
	if err := validateRawArgs(CmdFrom, rawArgs); err != nil {
		return nil, err
	}

	// named param mapping is done differently
	// because address for hosts contains char ':'
	var argMap map[string]string
	// if no pram named 'hosts' found, assume raw default
	if !strings.Contains(rawArgs, "hosts:") {
		rawArgs = makeNamedPram("hosts", rawArgs)
	}
	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("CAPTURE: %v", err)
	}

	cmd := &FromCommand{cmd: cmd{index: index, name: CmdFrom, args: argMap}}
	if err := validateCmdArgs(CmdFrom, argMap); err != nil {
		return nil, err
	}

	// populate machine representations
	for _, host := range spaceSep.Split(argMap["hosts"], -1) {
		addrVal := os.ExpandEnv(host)
		cmd.machines = append(cmd.machines, *NewMachine(addrVal))
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

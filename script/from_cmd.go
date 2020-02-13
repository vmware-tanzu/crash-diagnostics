// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Machine represents a machine as defined in FROM
type Machine struct {
	host,
	port,
	name string
}

// NewMachine returns a new *Machine
func NewMachine(host, port, name string) *Machine {
	if name == "" {
		name = fmt.Sprintf("%s:%s", host, port)
	}
	return &Machine{host: host, port: port, name: name}
}

// Address returns the host:port address
func (m *Machine) Address() string {
	return net.JoinHostPort(m.host, m.port)
}

// Host returns the host of the node address
func (m *Machine) Host() string {
	return m.host
}

// Port returns the port of the node address
func (m *Machine) Port() string {
	return m.port
}

// Name is a identifier for the machine
func (m *Machine) Name() string {
	return m.name
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

	// If params are not provided, assume nameless params format:
	// FROM addr.1 addr.2 addr.etc
	var argMap map[string]string
	if !strings.Contains(rawArgs, "hosts:") &&
		!strings.Contains(rawArgs, "nodes:") &&
		!strings.Contains(rawArgs, "source:") &&
		!strings.Contains(rawArgs, "port:") {
		rawArgs = makeNamedPram("hosts", rawArgs)
	}

	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("FROM: %v", err)
	}

	// add missing params
	if _, ok := argMap["port"]; !ok {
		argMap["port"] = Defaults.ServicePort
	}

	if _, ok := argMap["retries"]; !ok {
		argMap["retries"] = Defaults.ConnectionRetries
	}

	if _, ok := argMap["timeout"]; !ok {
		argMap["timeout"] = Defaults.ConnectionTimeout
	}

	if len(argMap["hosts"]) == 0 && len(argMap["nodes"]) == 0 {
		return nil, fmt.Errorf("FROM: must have hosts or nodes")
	}

	cmd := &FromCommand{cmd: cmd{index: index, name: CmdFrom, args: argMap}}
	if err := validateCmdArgs(CmdFrom, argMap); err != nil {
		return nil, err
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

// Hosts returns direct host address to source from
func (c *FromCommand) Hosts() []string {
	var result []string
	hosts := spaceSep.Split(strings.TrimSpace(c.cmd.args["hosts"]), -1)
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if len(host) == 0 {
			continue
		}
		result = append(result, ExpandEnv(host))
	}
	return result
}

// Nodes returns node names/ips
func (c *FromCommand) Nodes() []string {
	var result []string
	nodes := spaceSep.Split(c.cmd.args["nodes"], -1)
	for _, node := range nodes {
		node = strings.TrimSpace(node)
		if len(node) == 0 {
			continue
		}
		result = append(result, ExpandEnv(node))
	}
	return result
}

// Labels returns label filter used to select node to source from
func (c *FromCommand) Labels() string {
	return ExpandEnv(c.args[kubegetParams.labels])
}

// Port returns the default connection port
func (c *FromCommand) Port() string {
	return ExpandEnv(c.cmd.args["port"])
}

// ConnectionRetries returns the maximum number of connection retries
func (c *FromCommand) ConnectionRetries() int {
	str := ExpandEnv(c.cmd.args["retries"])
	val, err := strconv.Atoi(str)
	if err != nil {
		val = 30
	}
	return val
}

// ConnectionTimeout returns the duration to get a connection to servers
func (c *FromCommand) ConnectionTimeout() time.Duration {
	str := ExpandEnv(c.cmd.args["timeout"])
	to, err := time.ParseDuration(str)
	if err != nil {
		to = time.Second * 120
	}
	return to
}

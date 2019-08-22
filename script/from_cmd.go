package script

import (
	"fmt"
	"strings"
)

// Machine represents a source machine
type Machine struct {
	Address string
}

func (m Machine) String() string {
	return m.Address
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
		if arg != Defaults.FromValue {
			return nil, fmt.Errorf("%s only supports %s", CmdFrom, Defaults.FromValue)
		}
		cmd.machines = append(cmd.machines, Machine{Address: strings.TrimSpace(arg)})
		break // ignoring everything else. TODO fix.
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

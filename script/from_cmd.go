package script

import (
	"fmt"
	"strings"
)

// FromCommand represents instruction
// FROM <source>
// Where source can be:
// 1. Local : the current machine
// 2. List of machine names/addresses
// 3. cluster: uses Kuberentes cluster information to get list
type FromCommand struct {
	cmd
	sources []string
}

// NewFromCommand creates a value of type FromCommand
func NewFromCommand(index int, name string, args []string) (*FromCommand, error) {
	cmd := &FromCommand{cmd: cmd{index: index, name: name, args: args}}

	if err := validateCmdArgs(name, args); err != nil {
		return nil, err
	}

	for _, arg := range args {
		if arg != Defaults.FromValue {
			return nil, fmt.Errorf("%s only supports %s", CmdFrom, Defaults.FromValue)
		}
		cmd.sources = append(cmd.sources, strings.TrimSpace(arg))
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

func (c *FromCommand) Sources() []string {
	return c.sources
}

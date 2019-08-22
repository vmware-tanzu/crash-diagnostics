package script

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	envSep = regexp.MustCompile(`=`)
)

type EnvCommand struct {
	cmd
	envs []string
}

func NewEnvCommand(index int, args []string) (*EnvCommand, error) {
	cmd := &EnvCommand{cmd: cmd{index: index, name: CmdEnv, args: args}}

	if err := validateCmdArgs(CmdEnv, args); err != nil {
		return nil, err
	}

	for _, arg := range args {
		parts := envSep.Split(strings.TrimSpace(arg), -1)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid ENV arg %s", arg)
		}
		cmd.envs = append(cmd.envs, fmt.Sprintf("%s=%s", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])))
	}

	return cmd, nil
}

func (c *EnvCommand) Index() int {
	return c.cmd.index
}

func (c *EnvCommand) Name() string {
	return c.cmd.name
}

func (c *EnvCommand) Args() []string {
	return c.cmd.args
}

func (c *EnvCommand) Envs() []string {
	return c.envs
}

package script

// KubeConfigCommand represents Kubernetes configuration
type KubeConfigCommand struct {
	cmd
	kubeCfg string
}

// NewFromCommand creates a value of type FromCommand
func NewKubeConfigCommand(index int, args []string) (*KubeConfigCommand, error) {
	cmd := &KubeConfigCommand{cmd: cmd{index: index, name: CmdKubeConfig, args: args}}

	if err := validateCmdArgs(CmdKubeConfig, args); err != nil {
		return nil, err
	}
	cmd.kubeCfg = searchForConfig(args)
	return cmd, nil
}

func (c *KubeConfigCommand) Index() int {
	return c.cmd.index
}

func (c *KubeConfigCommand) Name() string {
	return c.cmd.name
}

func (c *KubeConfigCommand) Args() []string {
	return c.cmd.args
}

func (c *KubeConfigCommand) Config() string {
	return c.kubeCfg
}

// searchForConfig searches in several places for
// the kubernets config:
// 1. from passed args
// 2. from ENV variable
// 3. from local homedir
func searchForConfig(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return Defaults.KubeConfigValue
}

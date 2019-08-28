package script

import (
	"os"
	"path/filepath"
)

type CmdName string

var (
	CmdAs          = "AS"
	CmdCapture     = "CAPTURE"
	CmdCopy        = "COPY"
	CmdEnv         = "ENV"
	CmdFrom        = "FROM"
	CmdKubeConfig  = "KUBECONFIG"
	CmdFromDefault = "local"
	CmdWorkDir     = "WORKDIR"

	Defaults = struct {
		FromValue       string
		WorkdirValue    string
		KubeConfigValue string
	}{
		FromValue:    "local",
		WorkdirValue: "/tmp/flareout",
		KubeConfigValue: func() string {
			kubecfg := os.Getenv("KUBECONFIG")
			if kubecfg == "" {
				kubecfg = filepath.Join(os.Getenv("HOME"), ".kube", "config")
			}
			return kubecfg
		}(),
	}
)

type Script struct {
	Preambles map[string][]Command
	Actions   []Command
}

type CommandMeta struct {
	Name      string
	MinArgs   int
	MaxArgs   int
	Supported bool
}

var (
	Cmds = map[string]CommandMeta{
		CmdAs:         CommandMeta{Name: CmdAs, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdCapture:    CommandMeta{Name: CmdCapture, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdCopy:       CommandMeta{Name: CmdCopy, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdEnv:        CommandMeta{Name: CmdEnv, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdFrom:       CommandMeta{Name: CmdFrom, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdKubeConfig: CommandMeta{Name: CmdKubeConfig, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdWorkDir:    CommandMeta{Name: CmdWorkDir, MinArgs: 1, MaxArgs: 1, Supported: true},
	}
)

type Command interface {
	Index() int
	Name() string
	Args() []string
}

type cmd struct {
	index int
	name  string
	args  []string
}

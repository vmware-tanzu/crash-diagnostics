package script

type CmdName string

var (
	CmdAs          = "AS"
	CmdCapture     = "CAPTURE"
	CmdCopy        = "COPY"
	CmdEnv         = "ENV"
	CmdFrom        = "FROM"
	CmdFromDefault = "local"
	CmdWorkDir     = "WORKDIR"

	Defaults = struct {
		FromValue    string
		WorkdirValue string
	}{
		FromValue:    "local",
		WorkdirValue: "/tmp/flareout",
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
		CmdAs:      CommandMeta{Name: CmdAs, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdCapture: CommandMeta{Name: CmdCapture, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdCopy:    CommandMeta{Name: CmdCopy, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdEnv:     CommandMeta{Name: CmdEnv, MinArgs: 1, MaxArgs: -1, Supported: true},
		CmdFrom:    CommandMeta{Name: CmdFrom, MinArgs: 1, MaxArgs: 1, Supported: true},
		CmdWorkDir: CommandMeta{Name: CmdWorkDir, MinArgs: 1, MaxArgs: 1, Supported: true},
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

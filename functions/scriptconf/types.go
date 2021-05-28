package scriptconf

// Args represent input arguments passed to starlark function.
// Args can also be used as output arguments to built-in function.
//
// The argument map follows:
//   - error - used for output argument
//   - workdir string - a path that can be used as work directory during script exec
//   - gid string - the default group id to use when executing an OS command
//   - uid string - a default userid to use when executing an OS command
//   - default_shell string - path to a shell program that can be used as default (i.e. /bin/sh)
//   - requires [] string - a list of paths for commands that should be on the machine where script is executed
//   - use_ssh_agent bool - specifies if an ssh-agent should be setup for private key management
//
type Args struct {
	Workdir      string   `name:"workdir" optional:"true"`
	Gid          string   `name:"gid" optional:"true"`
	Uid          string   `name:"uid" optional:"true"`
	DefaultShell string   `name:"default_shell" optional:"true"`
	Requires     []string `name:"requires" optional:"true"`
	UseSSHAgent  bool     `name:"use_ssh_agent" optional:"true"`
}

// Result represent configuration returned by the function
// It is just an alias.
type Result struct {
	Error        string   `name:"error"`
	Workdir      string   `name:"workdir"`
	Gid          string   `name:"gid"`
	Uid          string   `name:"uid"`
	DefaultShell string   `name:"default_shell"`
	Requires     []string `name:"requires"`
	UseSSHAgent  bool     `name:"use_ssh_agent"`
}

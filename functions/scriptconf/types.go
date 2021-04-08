package scriptconf

// Params represent parameters for starlark function
type Params struct {
	Workdir      string
	Gid          string
	Uid          string
	DefaultShell string
	Requires     []string
	UseSSHAgent  bool
}

// Configuration represent configuration data returned by starlark function
type Configuration = Params

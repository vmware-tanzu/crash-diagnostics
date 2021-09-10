package scriptconf

import (
	"fmt"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/util"
	"go.starlark.net/starlark"
)

// ConfigFromThread retrieves script config result from provided
// thread instance. If found, bool = true.
func ConfigFromThread(t *starlark.Thread) (Config, bool) {
	if val := t.Local(Identifier); val != nil {
		result, ok := val.(Config)
		if !ok {
			return Config{}, ok
		}
		return result, true
	}
	return Config{}, false
}

// MakeDefaultConfigForThread creates a default config and adds it to the
// execution thread storage.
func MakeDefaultConfigForThread(t *starlark.Thread) (Config, error) {
	conf := makeDefaultConf()

	// create workdir if needed
	workdir, err := util.ExpandPath(conf.Workdir)
	if err != nil {
		return Config{}, err
	}
	if err := functions.MakeDir(workdir, 0744); err != nil && !os.IsExist(err) {
		return Config{}, fmt.Errorf("make thread config: failed to create workdir: %s", err)
	}

	// start local ssh-agent
	if conf.UseSSHAgent && t.Local(sshconf.AgentIdentifier) == nil {
		_, err := sshconf.MakeDefaultSSHAgentForThread(t)
		if err != nil {
			return Config{}, fmt.Errorf("make thread config: failed to start ssh agent: %s", err)
		}
	}

	// set conf in thread
	t.SetLocal(Identifier, conf)

	return conf, nil
}

func makeDefaultConf() Config {
	return Config{
		Workdir:      DefaultWorkdir(),
		Gid:          functions.DefaultGid(),
		Uid:          functions.DefaultUid(),
		DefaultShell: "/bin/sh",
		Requires:     []string{"/bin/ssh", "/bin/scp"},
		UseSSHAgent:  false,
	}
}

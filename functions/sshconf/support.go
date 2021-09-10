package sshconf

import (
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"go.starlark.net/starlark"
)

// SSHAgentFromThread retrieves ssh.Agent value from thread
func SSHAgentFromThread(t *starlark.Thread) (ssh.Agent, bool) {
	if agentVal := t.Local(AgentIdentifier); agentVal != nil {
		agent, ok := agentVal.(ssh.Agent)
		if !ok {
			return nil, false
		}
		return agent, true
	}
	return nil, false
}

// ConfigFromThread returns an sshconf.Config from provided
// starlark thread.
func ConfigFromThread(t *starlark.Thread) (Config, bool) {
	if confVal := t.Local(Identifier); confVal != nil {
		conf, ok := confVal.(Config)
		if !ok {
			return Config{}, false
		}
		return conf, true
	}
	return Config{}, false
}

// MakeDefaultConfigForThread creates a sshconf.Config value with
// default values, adds the config the the thread, and returns it.
func MakeDefaultConfigForThread(t *starlark.Thread) (Config, error) {
	conf := makeDefaultSSHConfig()

	// add private key to agent if agent was saved in thread
	if agent, ok := SSHAgentFromThread(t); ok {
		if err := agent.AddKey(conf.PrivateKeyPath); err != nil {
			return Config{}, fmt.Errorf("make thread config: unable to add private key to agent: %s", conf.PrivateKeyPath)
		}
	}

	return conf, nil
}

// MakeDefaultSSHAgentForThread starts an ssh agent and adds it to the
// specified thread.
func MakeDefaultSSHAgentForThread(t *starlark.Thread) (ssh.Agent, error) {
	agent, err := ssh.StartAgent()
	if err != nil {
		return nil, err
	}
	t.SetLocal(AgentIdentifier, agent)
	return agent, nil
}

func makeDefaultSSHConfig() Config {
	return Config{
		Username:       functions.DefaultUsername(),
		Port:           DefaultPort(),
		PrivateKeyPath: DefaultPKPath(),
		MaxRetries:     DefaultMaxRetries(),
		ConnTimeout:    DefaultConnTimeout(),
	}
}

package scriptconf

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		config Configuration
	}{
		{
			name:   "default values",
			params: Params{},
			config: Configuration{Workdir: defaultWorkdir, Gid: getGid(), Uid: getUid()},
		},
		{
			name:   "all values",
			params: Params{Workdir: "/a/b/c", Gid: "00", Uid: "01", UseSSHAgent: true, Requires: []string{"a/b"}},
			config: Configuration{Workdir: "/a/b/c", Gid: "00", Uid: "01", UseSSHAgent: true, Requires: []string{"a/b"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := Build(&starlark.Thread{}, test.params)
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Workdir != test.config.Workdir {
				t.Errorf("unexpected workdir value %s", cfg.Workdir)
			}
			if cfg.Gid != test.config.Gid {
				t.Errorf("unexpected Gid: %s", cfg.Gid)
			}
			if cfg.Uid != test.config.Uid {
				t.Errorf("unexpected Uid: %s", cfg.Uid)
			}
			if cfg.UseSSHAgent != test.config.UseSSHAgent {
				t.Errorf("unexpected UseSSHAgent: %t", cfg.UseSSHAgent)
			}
			if len(cfg.Requires) != len(test.config.Requires) {
				t.Errorf("unexpected len(Requires) %d", len(cfg.Requires))
			}
		})
	}
}

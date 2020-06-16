package starlark

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	identifiers = struct {
		crashdCfg string
		sshCfg    string
	}{
		crashdCfg: "crashd_config",
		sshCfg:    "ssh_config",
	}

	defaults = struct {
		crashdir    string
		workdir     string
		kubeconfig  string
		pkPath      string
		outPath     string
		connRetries int
		connTimeout int // seconds
	}{
		crashdir: filepath.Join(os.Getenv("HOME"), ".crashd"),
		workdir:  "/tmp/crashd",
		kubeconfig: func() string {
			kubecfg := os.Getenv("KUBECONFIG")
			if kubecfg == "" {
				kubecfg = filepath.Join(os.Getenv("HOME"), ".kube", "config")
			}
			return kubecfg
		}(),
		pkPath: func() string {
			return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
		}(),
		outPath:     "./crashd.tar.gz",
		connRetries: 10,
		connTimeout: 30,
	}
)

func isQuoted(val string) bool {
	single := `'`
	dbl := `"`
	if strings.HasPrefix(val, single) && strings.HasSuffix(val, single) {
		return true
	}
	if strings.HasPrefix(val, dbl) && strings.HasSuffix(val, dbl) {
		return true
	}
	return false
}

func trimQuotes(val string) string {
	unquoted, err := strconv.Unquote(val)
	if err != nil {
		return val
	}
	return unquoted
}

func getUsername() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Username
}

func getUid() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Uid
}

func getGid() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.Gid
}

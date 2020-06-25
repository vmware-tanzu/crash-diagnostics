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
		kubeCfg   string

		sshCfg         string
		port           string
		username       string
		privateKeyPath string
		maxRetries     string
		jumpUser       string
		jumpHost       string

		hostListProvider string
		hostResource     string
		resources        string
		run              string

		// Directives
		kubeCaptureDirective string
		kubeGetDirective     string
	}{
		crashdCfg: "crashd_config",
		kubeCfg:   "kube_config",

		sshCfg:         "ssh_config",
		port:           "port",
		username:       "username",
		privateKeyPath: "private_key_path",
		maxRetries:     "max_retries",
		jumpUser:       "jump_user",
		jumpHost:       "jump_host",

		hostListProvider: "host_list_provider",
		hostResource:     "host_resource",
		resources:        "resources",
		run:              "run",

		kubeGetDirective:     "kube_get",
		kubeCaptureDirective: "kube_capture",
	}

	defaults = struct {
		crashdir    string
		workdir     string
		kubeconfig  string
		sshPort     string
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
		sshPort: "22",
		pkPath: func() string {
			return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
		}(),
		outPath:     "./crashd.tar.gz",
		connRetries: 20,
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

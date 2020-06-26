package starlark

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var (
	strSanitization = regexp.MustCompile(`[^a-zA-Z0-9]`)

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
		capture          string

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
		capture:          "capture",

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

func getWorkdirFromThread(thread *starlark.Thread) (string, error) {
	val := thread.Local(identifiers.crashdCfg)
	if val == nil {
		return "", fmt.Errorf("%s not found in threard", identifiers.crashdCfg)
	}
	var result string
	if valStruct, ok := val.(*starlarkstruct.Struct); ok {
		if valStr, err := valStruct.Attr("workdir"); err == nil {
			if str, ok := valStr.(starlark.String); ok {
				result = string(str)
			}
		}
	}

	if len(result) == 0 {
		result = defaults.workdir
	}
	return result, nil
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

func captureOutput(source io.Reader, filePath, desc string) error {
	if source == nil {
		return fmt.Errorf("source reader is nill")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if len(desc) > 0 {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", desc)); err != nil {
			return err
		}
	}

	if _, err := io.Copy(file, source); err != nil {
		return err
	}

	logrus.Debugf("captured output in %s", filePath)

	return nil
}

func sanitizeStr(str string) string {
	return strSanitization.ReplaceAllString(str, "_")
}

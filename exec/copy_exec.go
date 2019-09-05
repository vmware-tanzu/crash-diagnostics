package exec

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/flare/script"
)

var (
	cliCpName = "cp"
	cliCpArgs = "-Rp"
)

func exeCopy(uid, gid int, dest string, cmd *script.CopyCommand) error {
	if _, err := exec.LookPath(cliCpName); err != nil {
		return err
	}

	for _, path := range cmd.Args() {
		if relPath, err := filepath.Rel(dest, path); err == nil && !strings.HasPrefix(relPath, "..") {
			logrus.Errorf("%s path %s cannot be relative to %s", cmd.Name(), path, dest)
			continue
		}

		logrus.Debugf("Copying %s to %s", path, dest)

		targetPath := filepath.Join(dest, path)
		targetDir := filepath.Dir(targetPath)
		if _, err := os.Stat(targetDir); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(targetDir, 0744); err != nil && !os.IsExist(err) {
					return err
				}
				logrus.Debugf("Created dir %s", targetDir)
			} else {
				return err
			}
		}

		args := []string{cliCpArgs, path, targetPath}
		_, err := CliRun(uint32(uid), uint32(gid), nil, cliCpName, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

package exec

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func exeCapture(asCmd *script.AsCommand, cmdCap *script.CaptureCommand, envs []string, workdir string) error {
	cmdStr := cmdCap.GetCliString()
	logrus.Debugf("Capturing CLI command %v", cmdStr)
	cliCmd, cliArgs := cmdCap.GetParsedCli()

	if _, err := exec.LookPath(cliCmd); err != nil {
		return err
	}

	asUid, asGid, err := asCmd.GetCredentials()
	if err != nil {
		return err
	}

	cmdReader, err := CliRun(uint32(asUid), uint32(asGid), envs, cliCmd, cliArgs...)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.txt", flatCmd(cmdStr))
	filePath := filepath.Join(workdir, fileName)
	logrus.Debugf("Capturing output of [%s] -into-> %s", cmdStr, filePath)
	if err := writeFile(cmdReader, filePath); err != nil {
		return err
	}

	return nil
}

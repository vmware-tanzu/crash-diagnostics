package exec

import (
	"fmt"
	"path/filepath"

	"gitlab.eng.vmware.com/vivienv/flare/ssh"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func exeRemotely(src *script.Script, machine *script.Machine, workdir string) error {
	asCmd, err := exeAs(src)
	if err != nil {
		return err
	}

	sshCmd, err := exeSSH(src)
	if err != nil {
		return err
	}

	user := asCmd.GetUserId()
	if sshCmd.GetUserId() != "" {
		user = sshCmd.GetUserId()
	}

	privKey := sshCmd.GetPrivateKeyPath()
	if privKey == "" {
		return fmt.Errorf("Missing private key file")
	}

	sshc := ssh.New(user, privKey)
	if err := sshc.Dial(machine.Address()); err != nil {
		return err
	}
	defer sshc.Hangup()

	for _, action := range src.Actions {
		switch cmd := action.(type) {
		case *script.CopyCommand:
			// if err := exeCopy(asCmd, cmd, workdir); err != nil {
			// 	return err
			// }
		case *script.CaptureCommand:
			// capture command output
			if err := captureRemotely(sshc, cmd, workdir); err != nil {
				return err
			}
		default:
			logrus.Errorf("Unsupported command %T", cmd)
		}
	}

	return nil
}

func captureRemotely(sshc *ssh.SSHClient, cmdCap *script.CaptureCommand, workdir string) error {
	cmdStr := cmdCap.GetCliString()
	logrus.Debugf("Capturing remote command command %v", cmdStr)
	cliCmd, cliArgs := cmdCap.GetParsedCli()

	cmdReader, err := sshc.SSHRun(cliCmd, cliArgs...)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.txt", sanitizeStr(cmdStr))
	filePath := filepath.Join(workdir, fileName)
	logrus.Debugf("Capturing output of [%s] -into-> %s", cmdStr, filePath)
	if err := writeFile(cmdReader, filePath); err != nil {
		return err
	}

	return nil
}

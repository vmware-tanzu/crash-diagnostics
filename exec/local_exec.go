// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

// exeLocally runs script using locally installed tool
func exeLocally(asCmd *script.AsCommand, action script.Command, workdir, output string) error {

	switch cmd := action.(type) {
	case *script.CopyCommand:
		if err := copyLocally(asCmd, cmd, workdir); err != nil {
			return err
		}
	case *script.CaptureCommand:
		// capture command output
		if err := captureLocally(asCmd, cmd, workdir, output); err != nil {
			return err
		}
	case *script.RunCommand:
		// run command and store result
		if err := runLocally(asCmd, cmd, output); err != nil {
			return err
		}
	default:
		logrus.Errorf("Unsupported command %T", cmd)
	}

	return nil
}

func captureLocally(asCmd *script.AsCommand, cmdCap *script.CaptureCommand, workdir, output string) error {
	cmdStr := cmdCap.GetCmdString()
	cliCmd, cliArgs, err := cmdCap.GetParsedCmd()
	if err != nil {
		return err
	}

	if _, err := exec.LookPath(cliCmd); err != nil {
		return err
	}

	asUid, asGid, err := asCmd.GetCredentials()
	if err != nil {
		return err
	}

	file, err := getFileForCaptureCmd(cmdStr, workdir, output)
	if err != nil {
		return err
	}

	// defer file close when not stdout/stderr
	switch output {
	case OutputStdout, OutputStderr:
	default:
		defer file.Close()
	}

	cmdReader, err := CliRun(uint32(asUid), uint32(asGid), cliCmd, cliArgs...)
	if err != nil {
		cliErr := fmt.Errorf("local command %s failed: %s", cliCmd, err)
		logrus.Warn(cliErr)
		return writeError(file, cliErr)
	}

	if err := writeFile(file, cmdReader); err != nil {
		return err
	}

	return nil
}

func runLocally(asCmd *script.AsCommand, cmdRun *script.RunCommand, output string) error {
	cmdStr := cmdRun.GetCmdString()
	cliCmd, cliArgs, err := cmdRun.GetParsedCmd()
	if err != nil {
		return err
	}

	if _, err := exec.LookPath(cliCmd); err != nil {
		return err
	}

	asUid, asGid, err := asCmd.GetCredentials()
	if err != nil {
		return err
	}

	logrus.Debugf("Running command [%s]", cmdStr)

	cmdReader, err := CliRun(uint32(asUid), uint32(asGid), cliCmd, cliArgs...)
	if err != nil {
		msgBytes, _ := ioutil.ReadAll(cmdReader)
		cmdErr := fmt.Errorf("RUN failed: command %s : %s : %s", cliCmd, err, strings.TrimSpace(string(msgBytes)))
		return cmdErr
	}

	cmdOutput := new(strings.Builder)
	if _, err := io.Copy(cmdOutput, cmdReader); err != nil {
		return fmt.Errorf("RUN failed: reading result of %s", err)
	}

	// save result of CMD
	result := strings.TrimSpace(string(cmdOutput.String()))
	if err := os.Setenv("CMD_RESULT", result); err != nil {
		return fmt.Errorf("RUN failed: setting CMD_RESULT: %s", err)
	}

	switch output {
	case OutputStdout:
		fmt.Fprintf(os.Stdout, "%s\n", result)
	case OutputStderr:
		fmt.Fprintf(os.Stderr, "%s\n", result)
	}

	return nil
}

var (
	cliCpShell      = "/bin/sh"
	cliCpShellParam = "-c"
	cliCpName       = "cp"
	cliCpArgs       = "-Rp"
)

func copyLocally(asCmd *script.AsCommand, cmd *script.CopyCommand, dest string) error {
	if _, err := exec.LookPath(cliCpName); err != nil {
		return err
	}

	asUid, asGid, err := asCmd.GetCredentials()
	if err != nil {
		return err
	}

	for _, path := range cmd.Paths() {
		if relPath, err := filepath.Rel(dest, path); err == nil && !strings.HasPrefix(relPath, "..") {
			logrus.Errorf("%s path %s cannot be relative to %s", cmd.Name(), path, dest)
			continue
		}

		// if path contains file pattern, adjust target
		pathDir, pathFile := filepath.Split(path)
		targetPath := filepath.Join(dest, path)
		targetDir := filepath.Dir(targetPath)
		if strings.Index(pathFile, "*") != -1 {
			targetPath = filepath.Join(dest, pathDir)
			targetDir = targetPath
		}

		if _, err := os.Stat(targetDir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if err := os.MkdirAll(targetDir, 0744); err != nil && !os.IsExist(err) {
				return err
			}
			logrus.Debugf("Created dir %s", targetDir)
		}

		logrus.Debugf("Copying %s to %s", path, targetPath)
		cpCmd := fmt.Sprintf("cp -Rp %s %s", path, targetPath)
		output, err := CliRun(uint32(asUid), uint32(asGid), cliCpShell, cliCpShellParam, cpCmd)
		if err != nil {
			msgBytes, _ := ioutil.ReadAll(output)
			cliErr := fmt.Errorf("local file copy failed: %s: %s: %s", path, string(msgBytes), err)
			logrus.Warn(cliErr)
		}
	}

	return nil
}

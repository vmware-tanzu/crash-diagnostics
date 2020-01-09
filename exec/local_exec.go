// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

// exeLocally runs script using locally installed tool
func exeLocally(asCmd *script.AsCommand, action script.Command, workdir string) error {
	switch cmd := action.(type) {
	case *script.CopyCommand:
		if err := copyLocally(asCmd, cmd, workdir); err != nil {
			return err
		}
	case *script.CaptureCommand:
		// capture command output
		if err := captureLocally(asCmd, cmd, nil, workdir); err != nil {
			return err
		}
	case *script.RunCommand:
		// run command and store result
		if err := runLocally(asCmd, cmd, workdir); err != nil {
			return err
		}
	default:
		logrus.Errorf("Unsupported command %T", cmd)
	}

	return nil
}

func captureLocally(asCmd *script.AsCommand, cmdCap *script.CaptureCommand, envs []string, workdir string) error {
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

	fileName := fmt.Sprintf("%s.txt", sanitizeStr(cmdStr))
	filePath := filepath.Join(workdir, fileName)
	logrus.Debugf("Capturing local command [%s] -into-> %s", cmdStr, filePath)

	cmdReader, err := CliRun(uint32(asUid), uint32(asGid), cliCmd, cliArgs...)
	if err != nil {
		cliErr := fmt.Errorf("local command %s failed: %s", cliCmd, err)
		logrus.Warn(cliErr)

		return writeCmdError(cliErr, filePath, cmdStr)
	}

	echo := false
	switch cmdCap.GetEcho() {
	case "true", "yes", "on":
		echo = true
	}

	if err := writeCmdOutput(cmdReader, filePath, echo, cmdStr); err != nil {
		return err
	}

	return nil
}

func runLocally(asCmd *script.AsCommand, cmdRun *script.RunCommand, workdir string) error {
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
		cmdErr := fmt.Errorf("Command failed: [%s]: %s", cliCmd, err)
		logrus.Error(cmdErr)
		return nil
	}

	bytes, err := ioutil.ReadAll(cmdReader)
	if err != nil {
		return fmt.Errorf("RUN: result: %s", err)
	}

	// save result of CMD
	result := strings.TrimSpace(string(bytes))
	if err := os.Setenv("CMD_RESULT", result); err != nil {
		return fmt.Errorf("RUN: set CMD_RESULT: %s", err)
	}

	switch cmdRun.GetEcho() {
	case "true", "yes", "on":
		fmt.Printf("%s\n%s\n", cmdRun.GetCmdString(), result)
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
		output, err := CliRun(uint32(asUid), uint32(asGid), "/bin/sh", "-c", cpCmd)
		if err != nil {
			msgBytes, _ := ioutil.ReadAll(output)
			cliErr := fmt.Errorf("local file copy failed: %s: %s: %s", path, string(msgBytes), err)
			logrus.Warn(cliErr)
			return writeCmdError(cliErr, targetPath, cpCmd)
		}
	}

	return nil
}

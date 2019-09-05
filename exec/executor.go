package exec

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/flare/script"
)

type Executor struct {
	script *script.Script
}

func New(src *script.Script) *Executor {
	return &Executor{script: src}
}

func (e *Executor) Execute() error {
	logrus.Info("Executing script file")
	// setup FROM
	fromCmds, ok := e.script.Preambles[script.CmdFrom]
	if !ok {
		return fmt.Errorf("Script missing valid %s", script.CmdFrom)
	}
	fromCmd := fromCmds[0].(*script.FromCommand)

	// setup AS instruction
	asCmds, ok := e.script.Preambles[script.CmdAs]
	if !ok {
		return fmt.Errorf("Script missing valid %s", script.CmdAs)
	}
	asCmd := asCmds[0].(*script.AsCommand)
	asUid, asGid, err := asCmd.GetCredentials()
	if err != nil {
		return err
	}
	logrus.Debugf("Executing as user %s:%s", asCmd.GetUserId(), asCmd.GetGroupId())

	// setup WORKDIR
	dirs, ok := e.script.Preambles[script.CmdWorkDir]
	if !ok {
		return fmt.Errorf("Script missing valid %s", script.CmdWorkDir)
	}
	workdir := dirs[0].(*script.WorkdirCommand)
	if err := os.MkdirAll(workdir.Dir(), 0744); err != nil && !os.IsExist(err) {
		return err
	}
	logrus.Debugf("Using workdir %s", workdir.Dir())

	// setup ENV
	var envPairs []string
	envCmds := e.script.Preambles[script.CmdEnv]
	for _, envCmd := range envCmds {
		env := envCmd.(*script.EnvCommand)
		if len(env.Envs()) > 0 {
			for _, arg := range env.Envs() {
				envPairs = append(envPairs, arg)
			}
		}
	}

	// retrieve KUBECONFIG and setup client connection
	cfgs, ok := e.script.Preambles[script.CmdKubeConfig]
	if !ok {
		return fmt.Errorf("Script missing valid %s", script.CmdKubeConfig)
	}
	cfgCmd := cfgs[0].(*script.KubeConfigCommand)
	if _, err := os.Stat(cfgCmd.Config()); err == nil {
		logrus.Debugf("Using KUBECONFIG %s", cfgCmd.Config())
		k8sClient, err := getK8sClient(cfgCmd.Config())
		if err != nil {
			logrus.Errorf("Failed to create Kubernetes API server client: %s", err)
		} else {
			if err := dumpClusterInfo(k8sClient, filepath.Join(workdir.Dir(), "cluster-dump.json")); err != nil {
				logrus.Errorf("Failed to retrieve cluster information: %s", err)
			}
		}
	} else {
		logrus.Warnf("Skipping cluster-info, unable to load KUBECONFIG %s: %s", cfgCmd.Config(), err)
	}

	// process actions for each cluster resource specified in FROM
	for _, fromMachine := range fromCmd.Machines() {
		machineAddr := fromMachine.Address
		if machineAddr != script.Defaults.FromValue {
			return fmt.Errorf("FROM only support 'local'")
		}
		machineWorkdir := filepath.Join(workdir.Dir(), machineAddr)
		if err := os.MkdirAll(machineWorkdir, 0744); err != nil && !os.IsExist(err) {
			return err
		}

		for _, action := range e.script.Actions {
			switch cmd := action.(type) {
			case *script.CopyCommand:
				if err := exeCopy(asUid, asGid, machineWorkdir, cmd); err != nil {
					return err
				}
			case *script.CaptureCommand:
				// capture command output
				cmdStr := cmd.GetCliString()
				logrus.Debugf("Capturing CLI command %v", cmdStr)
				cliCmd, cliArgs := cmd.GetParsedCli()

				if _, err := exec.LookPath(cliCmd); err != nil {
					return err
				}

				cmdReader, err := CliRun(uint32(asUid), uint32(asGid), envPairs, cliCmd, cliArgs...)
				if err != nil {
					return err
				}

				fileName := fmt.Sprintf("%s.txt", flatCmd(cmdStr))
				filePath := filepath.Join(machineWorkdir, fileName)
				logrus.Debugf("Capturing output of [%s] -into-> %s", cmdStr, filePath)
				if err := writeFile(cmdReader, filePath); err != nil {
					return err
				}
			default:
			}
		}
	}
	return nil
}

func writeFile(source io.Reader, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, source); err != nil {
		return err
	}
	logrus.Debugf("Wrote file %s", filePath)

	return nil
}

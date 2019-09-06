package exec

import (
	"fmt"
	"io"
	"os"
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
	// exec FROM
	fromCmd, err := exeFrom(e.script)
	if err != nil {
		return err
	}

	// exec WORKDIR
	workdir, err := exeWorkdir(e.script)
	if err != nil {
		return nil
	}
	logrus.Debugf("Using workdir %s", workdir.Dir())

	// retrieve KUBECONFIG and setup client connection
	exeClusterInfo(e.script, filepath.Join(workdir.Dir(), "cluster-dump.json"))

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

		switch machineAddr {
		case "local":
			if err := exeLocally(e.script, machineWorkdir); err != nil {
				return err
			}
		default:
			// exeRemotely(e.script, machineWorkdir)
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

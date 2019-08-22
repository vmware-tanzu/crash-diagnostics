package exec

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	// process action for each FROM source

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
				// TODO - COPY uses a go implementation which means uid/guid
				// for the COPY cmd cannot be applied using the flare file.
				// This may need to be changed to a os/cmd external call

				// walk each arg and copy to workdir
				for _, path := range cmd.Args() {
					if relPath, err := filepath.Rel(machineWorkdir, path); err == nil && !strings.HasPrefix(relPath, "..") {
						logrus.Errorf("%s path %s cannot be relative to workdir %s", cmd.Name(), path, machineWorkdir)
						continue
					}
					logrus.Debugf("Copying content from %s", path)

					err := filepath.Walk(path, func(file string, finfo os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						relPath := file
						if filepath.IsAbs(file) {
							relPath, err = filepath.Rel("/", file)
							if err != nil {
								return err
							}
						}

						// setup subpath where source is copied to
						subpath := filepath.Join(machineWorkdir, relPath)
						subpathDir := filepath.Dir(subpath)
						if _, err := os.Stat(subpathDir); err != nil && os.IsNotExist(err) {
							if err := os.MkdirAll(subpathDir, 0744); err != nil && !os.IsExist(err) {
								return err
							}
							logrus.Debugf("Created parent dir %s", subpathDir)
						}

						switch {
						case finfo.Mode().IsDir():
							if err := os.MkdirAll(subpath, 0744); err != nil && !os.IsExist(err) {
								return err
							}
							logrus.Debugf("Created subpath %s", subpath)
							return nil
						case finfo.Mode().IsRegular():
							logrus.Debugf("Copying %s -> %s", file, subpath)
							srcFile, err := os.Open(file)
							if err != nil {
								return err
							}
							defer srcFile.Close()

							desFile, err := os.Create(subpath)
							if err != nil {
								return err
							}
							n, err := io.Copy(desFile, srcFile)
							if closeErr := desFile.Close(); closeErr != nil {
								return closeErr
							}
							if err != nil {
								return err
							}

							if n != finfo.Size() {
								return fmt.Errorf("%s did not complet for %s", cmd.Name(), file)
							}
						default:
							return fmt.Errorf("%s unknown file type for %s", cmd.Name(), file)
						}
						return nil
					})

					if err != nil {
						logrus.Error(err)
					}
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

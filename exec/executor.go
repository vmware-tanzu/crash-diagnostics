package exec

import (
	"fmt"
	"io"
	"os"
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
	logrus.Info("Executing flare file")
	// setup FROM
	fromArg := script.Defaults.FromValue
	from := e.script.Preambles[script.CmdFrom]
	if from != nil && len(from.Args) > 0 {
		fromArg = from.Args[0]
	}
	if fromArg != script.Defaults.FromValue {
		return fmt.Errorf("%s only supports %s", script.CmdFrom, script.Defaults.FromValue)
	}
	logrus.Debugf("Collecting data from machine %s", fromArg)

	// TODO setup guid / uid

	// setup WORKDIR
	workdir := "/tmp/flareout"
	dir := e.script.Preambles[script.CmdWorkDir]
	if dir != nil && len(dir.Args) > 0 {
		workdir = dir.Args[0]
	}
	if err := os.MkdirAll(workdir, 0744); err != nil && !os.IsExist(err) {
		return err
	}
	logrus.Debugf("Using workdir %s", workdir)

	// process actions
	for _, cmd := range e.script.Actions {
		switch cmd.Name {
		case script.CmdCopy:
			if len(cmd.Args) < script.Cmds[script.CmdCopy].MinArgs {
				logrus.Errorf("%s missing argument, skipping", cmd.Name)
				continue
			}
			// walk each arg and copy to workdir
			for _, path := range cmd.Args {
				if relPath, err := filepath.Rel(workdir, path); err == nil && !strings.HasPrefix(relPath, "..") {
					logrus.Errorf("%s path %s cannot be relative to workdir %s", cmd.Name, path, workdir)
					continue
				}
				logrus.Debugf("Copying content from %s", path)

				err := filepath.Walk(path, func(file string, finfo os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					//TODO subpath calculation flattens the file source, that's wrong.
					// subpath should include full path of file, not just the base.
					subpath := filepath.Join(workdir, filepath.Base(file))
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
							return fmt.Errorf("%s did not complet for %s", cmd.Name, file)
						}
					default:
						return fmt.Errorf("%s unknown file type for %s", cmd.Name, file)
					}
					return nil
				})

				if err != nil {
					logrus.Error(err)
				}
			}
		case script.CmdCapture:
			if len(cmd.Args) < script.Cmds[script.CmdCopy].MinArgs {
				logrus.Errorf("%s missing argument", cmd.Name)
				continue
			}

			// capture command output
			cmdStr := strings.Join(cmd.Args, " ")
			logrus.Debugf("Parsing CLI command %v", cmdStr)
			cliCmd, cliArgs := CliParse(cmdStr)
			if cliCmd == "" {
				logrus.Debug("Skipping empty command")
				continue
			}
			cmdReader, err := CliRun(cliCmd, cliArgs...)
			if err != nil {
				return err
			}
			fileName := fmt.Sprintf("%s.txt", flatCmd(cmdStr))
			filePath := filepath.Join(workdir, fileName)
			logrus.Debugf("Capturing command out: [%s] -> %s", cmdStr, filePath)
			if err := writeFile(cmdReader, filePath); err != nil {
				return err
			}

		default:
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
	return nil
}

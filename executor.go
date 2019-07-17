package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type executor struct {
	script *script
}

func newExecutor(src *script) *executor {
	return &executor{script: src}
}

func (e *executor) exec() error {
	// setup FROM
	fromArg := cmdFromDefault
	from := e.script.preambles[cmdFrom]
	if from != nil && len(from.args) > 0 {
		fromArg = from.args[0]
	}
	if fromArg != cmdFromDefault {
		return fmt.Errorf("%s only supports %s", cmdFrom, cmdFromDefault)
	}

	// TODO setup guid / uid

	// setup WORKDIR
	workdir := "/tmp/flareout"
	dir := e.script.preambles[cmdWorkDir]
	if dir != nil && len(dir.args) > 0 {
		workdir = dir.args[0]
	}
	if err := os.MkdirAll(workdir, 0744); err != nil && !os.IsExist(err) {
		return err
	}

	// process actions
	for _, cmd := range e.script.actions {
		switch cmd.name {
		case cmdCopy:
			if len(cmd.args) < cmds[cmdCopy].minArgs {
				return fmt.Errorf("%s missing argument", cmd.name)
			}
			// walk each arg and copy to workdir
			for _, path := range cmd.args {
				if relPath, err := filepath.Rel(workdir, path); err == nil && !strings.HasPrefix(relPath, "..") {
					return fmt.Errorf("%s path %s cannot be relative to workdir %s", cmd.name, path, workdir)
				}

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
						return nil
					case finfo.Mode().IsRegular():
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
							return fmt.Errorf("%s did not complet for %s", cmd.name, file)
						}
					default:
						return fmt.Errorf("%s unknown file type for %s", cmd.name, file)
					}

					return nil
				})
				if err != nil {
					return err
				}
			}
		case cmdCapture:
		default:
		}
	}

	return nil
}

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
)

type Executor interface {
	Execute() (io.Reader, error)
}

type CommandExecutor struct {
	command *exec.Cmd
	output  *bytes.Buffer
}

func NewCommandExecutor(cmd string, args ...string) *CommandExecutor {
	executor := new(CommandExecutor)
	executor.output = new(bytes.Buffer)

	executor.command = exec.Command(cmd, args...)
	executor.command.Stdout = executor.output
	executor.command.Stderr = executor.output

	return executor
}

func (e *CommandExecutor) Execute() (io.Reader, error) {
	if err := e.runCmd(); err != nil {
		return nil, err
	}
	return e.output, nil
}

func (e *CommandExecutor) ExecToFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
	}()

	reader, err := e.Execute()
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

func (e *CommandExecutor) runCmd() error {
	if e.command == nil {
		return errors.New("missing command")
	}
	if err := e.command.Run(); err != nil {
		return err
	}
	return nil
}

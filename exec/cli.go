package exec

import (
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"syscall"
)

var (
	spaceSep = regexp.MustCompile(`\s`)
	quoteSet = regexp.MustCompile(`[\"\']`)
)

func CliRun(cmd string, args ...string) (io.Reader, error) {
	command, output := prepareCmd(cmd, args...)

	if err := command.Run(); err != nil {
		return nil, err
	}

	return output, nil
}

func CliRunAs(uid, gid uint32, cmd string, args ...string) (io.Reader, error) {
	command, output := prepareCmd(cmd, args...)
	command.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}

	if err := command.Run(); err != nil {
		return nil, err
	}

	return output, nil
}

func CliParse(cmdStr string) (cmd string, args []string) {
	args = []string{}
	parts := spaceSep.Split(cmdStr, -1)
	if len(parts) == 0 {
		return
	}
	if len(parts) == 1 {
		cmd = parts[0]
		return
	}
	cmd = parts[0]
	args = parts[1:]
	return
}

func prepareCmd(cmd string, args ...string) (*exec.Cmd, io.Reader) {
	output := new(bytes.Buffer)
	command := exec.Command(cmd, args...)
	command.Stdout = output
	command.Stderr = output
	return command, output
}

func flatCmd(cmd string) string {
	str := quoteSet.ReplaceAllString(cmd, "")
	return spaceSep.ReplaceAllString(str, "_")
}

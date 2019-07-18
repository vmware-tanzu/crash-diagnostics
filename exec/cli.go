package exec

import (
	"bytes"
	"io"
	"os/exec"
	"regexp"
)

var (
	spaceSep = regexp.MustCompile(`\s`)
	quoteSet = regexp.MustCompile(`[\"\']`)
)

func CliRun(cmd string, args ...string) (io.Reader, error) {
	output := new(bytes.Buffer)

	command := exec.Command(cmd, args...)
	command.Stdout = output
	command.Stderr = output

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

func flatCmd(cmd string) string {
	str := quoteSet.ReplaceAllString(cmd, "")
	return spaceSep.ReplaceAllString(str, "_")
}

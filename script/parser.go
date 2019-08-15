package script

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var cmdSep = regexp.MustCompile(`\s`)

func Parse(reader io.Reader) (*Script, error) {
	logrus.Debugln("Parsing flare script")
	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	var script Script
	script.Preambles = make(map[string][]Command)
	line := 1
	for lineScanner.Scan() {
		text := strings.TrimSpace(lineScanner.Text())
		if text == "" {
			line++
			continue
		}
		logrus.Debugf("Parsing [%d: %s]", line, text)
		tokens := cmdSep.Split(text, -1)
		cmdName := tokens[0]
		if !Cmds[cmdName].Supported {
			return nil, fmt.Errorf("line %d: %s unsupported", line, cmdName)
		}
		// TODO additional validation needed:
		// 1) validate preambles and args
		// 2) validate each action and args
		switch cmdName {
		case CmdAs, CmdFrom, CmdWorkDir, CmdEnv:
			cmd := Command{Index: line, Name: cmdName, Args: tokens[1:]}
			script.Preambles[cmdName] = append(script.Preambles[cmdName], cmd)
		case CmdCopy:
			cmd := Command{Index: line, Name: cmdName, Args: tokens[1:]}
			script.Actions = append(script.Actions, cmd)
		case CmdCapture:
			cmdStr := strings.Join(tokens[1:], " ")
			command := Command{Index: line, Name: cmdName, Args: []string{cmdStr}}
			script.Actions = append(script.Actions, command)
		default:
			return nil, fmt.Errorf("%s not supported", cmdName)
		}
		logrus.Debugf("%s parsed OK", cmdName)
		line++
	}
	logrus.Debug("Done parsing")
	return &script, nil
}

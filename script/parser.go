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
	script.Preambles = make(map[string]*Command)
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
		cmd := Command{Index: line, Name: cmdName, Args: tokens[1:]}
		switch cmdName {
		case CmdAs, CmdFrom, CmdWorkDir:
			logrus.Debugf("Preamble encountered: %s", cmdName)
			script.Preambles[cmdName] = &cmd
		case CmdCopy:
			logrus.Debug("COPY action encountered")
			script.Actions = append(script.Actions, cmd)
		case CmdCapture:
			logrus.Debug("CAPTURE action encountered")
			script.Actions = append(script.Actions, cmd)
		default:
			return nil, fmt.Errorf("%s not supported", cmdName)
		}
		line++
	}
	logrus.Debug("Done parsing")
	return &script, nil
}

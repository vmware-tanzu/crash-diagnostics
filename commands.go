package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type cmdMeta struct {
	name      string
	minArgs   int
	supported bool
}

type command struct {
	index int
	name  string
	args  []string
}

type script struct {
	preambles map[string]*command
	actions   []command
}

var (
	cmdSep = regexp.MustCompile(`\s`)

	cmdAs          = "AS"
	cmdCapture     = "CAPTURE"
	cmdCopy        = "COPY"
	cmdEnv         = "ENV"
	cmdFrom        = "FROM"
	cmdFromDefault = "local"
	cmdWorkDir     = "WORKDIR"

	cmds = map[string]cmdMeta{
		cmdAs:      cmdMeta{name: cmdAs, minArgs: 1, supported: true},
		cmdCapture: cmdMeta{name: cmdCapture, minArgs: 1, supported: true},
		cmdCopy:    cmdMeta{name: cmdCopy, minArgs: 1, supported: true},
		cmdEnv:     cmdMeta{name: cmdEnv, minArgs: 1, supported: true},
		cmdFrom:    cmdMeta{name: cmdFrom, minArgs: 1, supported: true},
		cmdWorkDir: cmdMeta{name: cmdWorkDir, minArgs: 1, supported: true},
	}
)

func parse(reader io.Reader) (*script, error) {
	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	var script script
	script.preambles = make(map[string]*command)
	line := 1
	for lineScanner.Scan() {
		tokens := cmdSep.Split(lineScanner.Text(), -1)
		cmdName := tokens[0]
		if !cmds[cmdName].supported {
			return nil, fmt.Errorf("line %d: %s unsupported", line, cmdName)
		}
		// TODO additional validation needed:
		// 1) validate preambles and args
		// 2) validate each action and args
		cmd := command{index: line, name: cmdName, args: tokens[1:]}
		switch cmdName {
		case cmdAs, cmdFrom, cmdWorkDir:
			script.preambles[cmdName] = &cmd
		case cmdCopy:
			script.actions = append(script.actions, cmd)
		case cmdCapture:
			script.actions = append(script.actions, cmd)
		default:
			return nil, fmt.Errorf("%s not supported", cmdName)
		}
		line++
	}
	return &script, nil
}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/script"
)

var (
	spaceSep       = regexp.MustCompile(`\s`)
)

// Parse parses the textual script into an *script.Script representation
func Parse(reader io.Reader) (*script.Script, error) {
	logrus.Info("Parsing scr file")

	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	var scr script.Script
	scr.Preambles = make(map[string][]script.Command)
	line := 1
	for lineScanner.Scan() {
		text := strings.TrimSpace(lineScanner.Text())
		if text == "" || text[0] == '#' {
			line++
			continue
		}
		logrus.Debugf("Parsing [%d: %s]", line, text)

		// split DIRECTIVE [ARGS] after first space(s)
		var cmdName, rawArgs string
		tokens := spaceSep.Split(text, 2)
		if len(tokens) == 2 {
			rawArgs = tokens[1]
		}
		cmdName = tokens[0]

		if !script.Cmds[cmdName].Supported {
			return nil, fmt.Errorf("line %d: %s unsupported", line, cmdName)
		}

		switch cmdName {
		case script.CmdAs:
			cmd, err := script.NewAsCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Preambles[script.CmdAs] = []script.Command{cmd} // save only last AS instruction
		case script.CmdEnv:
			cmd, err := script.NewEnvCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Preambles[script.CmdEnv] = append(scr.Preambles[script.CmdEnv], cmd)
		case script.CmdFrom:
			cmd, err := script.NewFromCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Preambles[script.CmdFrom] = []script.Command{cmd} // saves only last FROM
		case script.CmdKubeConfig:
			cmd, err := script.NewKubeConfigCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Preambles[script.CmdKubeConfig] = []script.Command{cmd}
		case script.CmdAuthConfig:
			cmd, err := script.NewAuthConfigCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Preambles[script.CmdAuthConfig] = []script.Command{cmd}
		case script.CmdOutput:
			cmd, err := script.NewOutputCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Preambles[script.CmdOutput] = []script.Command{cmd}
		case script.CmdWorkDir:
			cmd, err := script.NewWorkdirCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Preambles[script.CmdWorkDir] = []script.Command{cmd}
		case script.CmdCapture:
			cmd, err := script.NewCaptureCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Actions = append(scr.Actions, cmd)
		case script.CmdCopy:
			cmd, err := script.NewCopyCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Actions = append(scr.Actions, cmd)
		case script.CmdRun:
			cmd, err := script.NewRunCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Actions = append(scr.Actions, cmd)
		case script.CmdKubeGet:
			cmd, err := script.NewKubeGetCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			scr.Actions = append(scr.Actions, cmd)
		default:
			return nil, fmt.Errorf("%s not supported", cmdName)
		}
		logrus.Debugf("%s parsed OK", cmdName)
		line++
	}
	logrus.Debug("Done parsing")
	return enforceDefaults(&scr)
}

// enforceDefaults adds missing defaults to the script
func enforceDefaults(scr *script.Script) (*script.Script, error) {
	logrus.Debug("Applying default values")
	if _, ok := scr.Preambles[script.CmdAs]; !ok {
		cmd, err := script.NewAsCommand(0, fmt.Sprintf("userid:%d groupid:%d", os.Getuid(), os.Getgid()))
		if err != nil {
			return scr, err
		}
		logrus.Debugf("AS %s:%s (as default)", cmd.GetUserId(), cmd.GetGroupId())
		scr.Preambles[script.CmdAs] = []script.Command{cmd}
	}

	if _, ok := scr.Preambles[script.CmdFrom]; !ok {
		cmd, err := script.NewFromCommand(0, script.Defaults.FromValue)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("FROM %v (as default)", cmd.Nodes())
		scr.Preambles[script.CmdFrom] = []script.Command{cmd}
	}
	if _, ok := scr.Preambles[script.CmdAuthConfig]; !ok {
		cmd, err := script.NewAuthConfigCommand(0, fmt.Sprintf("username:${USER} private-key:${HOME}/.ssh/id_rsa"))
		if err != nil {
			return nil, err
		}
		logrus.Debug("AUTHCONFIG set with default")
		scr.Preambles[script.CmdAuthConfig] = []script.Command{cmd}
	}
	if _, ok := scr.Preambles[script.CmdWorkDir]; !ok {
		cmd, err := script.NewWorkdirCommand(0, fmt.Sprintf("path:%s", script.Defaults.WorkdirValue))
		if err != nil {
			return nil, err
		}
		logrus.Debugf("WORKDIR %s (as default)", cmd.Path())
		scr.Preambles[script.CmdWorkDir] = []script.Command{cmd}
	}

	if _, ok := scr.Preambles[script.CmdOutput]; !ok {
		cmd, err := script.NewOutputCommand(0, fmt.Sprintf("path:%s", script.Defaults.OutputValue))
		if err != nil {
			return nil, err
		}
		logrus.Debugf("OUTPUT %s (as default)", cmd.Path())
		scr.Preambles[script.CmdOutput] = []script.Command{cmd}
	}

	if _, ok := scr.Preambles[script.CmdKubeConfig]; !ok {
		cmd, err := script.NewKubeConfigCommand(0, fmt.Sprintf("path:%s", script.Defaults.KubeConfigValue))
		if err != nil {
			return nil, err
		}
		logrus.Debugf("KUBECONFIG %s (as default)", cmd.Path())
		scr.Preambles[script.CmdKubeConfig] = []script.Command{cmd}
	}
	return scr, nil
}

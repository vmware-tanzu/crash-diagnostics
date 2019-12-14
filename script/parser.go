// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	spaceSep = regexp.MustCompile(`\s`)
	paramSep = regexp.MustCompile(`:`)
	quoteSet = regexp.MustCompile(`[\"\']`)
	cmdSep   = regexp.MustCompile(`\s`)
)

// Parse parses the textual script from reader into an *Script representation
func Parse(reader io.Reader) (*Script, error) {
	logrus.Info("Parsing script file")

	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	var script Script
	script.Preambles = make(map[string][]Command)
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

		if !Cmds[cmdName].Supported {
			return nil, fmt.Errorf("line %d: %s unsupported", line, cmdName)
		}

		switch cmdName {
		case CmdAs:
			cmd, err := NewAsCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Preambles[CmdAs] = []Command{cmd} // save only last AS instruction
		case CmdEnv:
			cmd, err := NewEnvCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Preambles[CmdEnv] = append(script.Preambles[CmdEnv], cmd)
		case CmdFrom:
			cmd, err := NewFromCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Preambles[CmdFrom] = []Command{cmd} // saves only last FROM
		case CmdKubeConfig:
			cmd, err := NewKubeConfigCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Preambles[CmdKubeConfig] = []Command{cmd}
		case CmdAuthConfig:
			cmd, err := NewAuthConfigCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Preambles[CmdAuthConfig] = []Command{cmd}
		case CmdOutput:
			cmd, err := NewOutputCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Preambles[CmdOutput] = []Command{cmd}
		case CmdWorkDir:
			cmd, err := NewWorkdirCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Preambles[CmdWorkDir] = []Command{cmd}
		case CmdCapture:
			cmd, err := NewCaptureCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Actions = append(script.Actions, cmd)
		case CmdCopy:
			cmd, err := NewCopyCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Actions = append(script.Actions, cmd)
		case CmdRun:
			cmd, err := NewRunCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Actions = append(script.Actions, cmd)
		case CmdKubeGet:
			cmd, err := NewKubeGetCommand(line, rawArgs)
			if err != nil {
				return nil, err
			}
			script.Actions = append(script.Actions, cmd)
		default:
			return nil, fmt.Errorf("%s not supported", cmdName)
		}
		logrus.Debugf("%s parsed OK", cmdName)
		line++
	}
	logrus.Debug("Done parsing")
	return enforceDefaults(&script)
}

func validateRawArgs(cmdName, rawArgs string) error {
	cmd, ok := Cmds[cmdName]
	if !ok {
		return fmt.Errorf("%s is unknown", cmdName)
	}
	if len(rawArgs) == 0 && cmd.MinArgs > 0 {
		return fmt.Errorf("%s must have at least %d argument(s)", cmdName, cmd.MinArgs)
	}
	return nil
}

func validateCmdArgs(cmdName string, args map[string]string) error {
	cmd, ok := Cmds[cmdName]
	if !ok {
		return fmt.Errorf("%s is unknown", cmdName)
	}

	minArgs := cmd.MinArgs
	maxArgs := cmd.MaxArgs

	if len(args) < minArgs {
		return fmt.Errorf("%s must have at least %d argument(s)", cmdName, minArgs)
	}

	if maxArgs > -1 && len(args) > maxArgs {
		return fmt.Errorf("%s can only have up to %d argument(s)", cmdName, maxArgs)
	}

	return nil
}

// mapArgs takes the rawArgs in the form of
//
//    param0:"val0" param1:"val1" ... paramN:"valN"
//
// The param name must be followed by a colon and the value
// may be quoted or unquoted. It is an error if
// split(rawArgs[n], ":") yields to a len(slice) < 2.
func mapArgs(rawArgs string) (map[string]string, error) {
	argMap := make(map[string]string)

	// split params: param0:<param-val0> paramN:<param-valN> badparam
	params, err := wordSplit(rawArgs)
	if err != nil {
		return nil, err
	}

	// for each, split pram:<pram-value> into {param, <param-val>}
	for _, param := range params {
		parts := paramSep.Split(param, 2)
		if len(parts) != 2 {
			return argMap, fmt.Errorf("invalid param: %s", param)
		}
		name := parts[0]
		val := trimQuotes(parts[1])
		argMap[name] = val
	}

	return argMap, nil
}

// isNamedParam returs true if str has the form
//
//    name:value
//
func isNamedParam(str string) bool {
	if len(str) == 0 {
		return false
	}

	parts := paramSep.Split(str, 2)
	if len(parts) >= 2 {
		return true
	}
	return false
}

// makeParam
func makeNamedPram(name, value string) string {
	value = strings.TrimSpace(value)
	// possibly already quoted
	if value[0] == '"' || value[0] == '\'' {
		return fmt.Sprintf("%s:%s", name, value)
	}
	// return as quoted
	return fmt.Sprintf(`%s:'%s'`, name, value)
}

// enforceDefaults adds missing defaults to the script
func enforceDefaults(script *Script) (*Script, error) {
	logrus.Debug("Applying default values")
	if _, ok := script.Preambles[CmdAs]; !ok {
		cmd, err := NewAsCommand(0, fmt.Sprintf("userid:%d groupid:%d", os.Getuid(), os.Getgid()))
		if err != nil {
			return script, err
		}
		logrus.Debugf("AS %s:%s (as default)", cmd.GetUserId(), cmd.GetGroupId())
		script.Preambles[CmdAs] = []Command{cmd}
	}

	if _, ok := script.Preambles[CmdFrom]; !ok {
		cmd, err := NewFromCommand(0, Defaults.FromValue)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("FROM %v (as default)", cmd.Machines())
		script.Preambles[CmdFrom] = []Command{cmd}
	}

	if _, ok := script.Preambles[CmdWorkDir]; !ok {
		cmd, err := NewWorkdirCommand(0, fmt.Sprintf("path:%s", Defaults.WorkdirValue))
		if err != nil {
			return nil, err
		}
		logrus.Debugf("WORKDIR %s (as default)", cmd.Path())
		script.Preambles[CmdWorkDir] = []Command{cmd}
	}

	if _, ok := script.Preambles[CmdOutput]; !ok {
		cmd, err := NewOutputCommand(0, fmt.Sprintf("path:%s", Defaults.OutputValue))
		if err != nil {
			return nil, err
		}
		logrus.Debugf("OUTPUT %s (as default)", cmd.Path())
		script.Preambles[CmdOutput] = []Command{cmd}
	}

	if _, ok := script.Preambles[CmdKubeConfig]; !ok {
		cmd, err := NewKubeConfigCommand(0, fmt.Sprintf("path:%s", Defaults.KubeConfigValue))
		if err != nil {
			return nil, err
		}
		logrus.Debugf("KUBECONFIG %s (as default)", cmd.Path())
		script.Preambles[CmdKubeConfig] = []Command{cmd}
	}
	return script, nil
}

func cmdParse(cmdStr string) (cmd string, args []string, err error) {
	logrus.Debugf("Parsing: %s", cmdStr)
	args, err = wordSplit(cmdStr)
	if err != nil {
		return "", nil, err
	}
	return args[0], args[1:], nil
}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/vmware-tanzu/crash-diagnostics/script"
)

var (
	spaceSep       = regexp.MustCompile(`\s`)
	namedParamRegx = regexp.MustCompile(`^([a-z0-9_\-]+)(:)(["']{0,1}.+["']{0,1})$`)
	directiveRegx  = regexp.MustCompile(`^([A-Z0-9_\-]+)\s+(.+)$`)
)

// Parse parses the textual script into an *script.Script representation
func Parse(reader io.Reader) (*script.Script, error) {
	logrus.Info("Parsing scr file")

	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	scr := script.New()
	line := 1

	for lineScanner.Scan() {
		text := strings.TrimSpace(lineScanner.Text())
		if text == "" || text[0] == '#' {
			line++
			continue
		}
		logrus.Debugf("Parsing [%d: %s]", line, text)

		// split DIRECTIVE [ARGS] after first space(s)
		//var cmdName, rawArgs string
		//tokens := spaceSep.Split(text, 2)
		//if len(tokens) == 2 {
		//	rawArgs = tokens[1]
		//}
		//cmdName = tokens[0]

		cmdName, rawArgs, err := splitDirectiveLine(text)
		if err != nil {
			return nil, fmt.Errorf("invalid directive: %s", err)
		}

		switch cmdName {
		case script.CmdAs, script.CmdAsConfig:
			scr.AddConfigDirective(line, script.CmdAsConfig, rawArgs)

		//case script.CmdEnv:
		//	cmd, err := script.NewEnvCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Preambles[script.CmdEnv] = append(scr.Preambles[script.CmdEnv], cmd)
		//case script.CmdFrom:
		//	cmd, err := script.NewFromCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Preambles[script.CmdFrom] = []script.Directive{cmd} // saves only last FROM
		//case script.CmdKubeConfig:
		//	cmd, err := script.NewKubeConfigCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Preambles[script.CmdKubeConfig] = []script.Directive{cmd}
		//case script.CmdAuthConfig:
		//	cmd, err := script.NewAuthConfigCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Preambles[script.CmdAuthConfig] = []script.Directive{cmd}
		//case script.CmdOutput:
		//	cmd, err := script.NewOutputCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Preambles[script.CmdOutput] = []script.Directive{cmd}
		//case script.CmdWorkDir:
		//	cmd, err := script.NewWorkdirCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Preambles[script.CmdWorkDir] = []script.Directive{cmd}
		//case script.CmdCapture:
		//	cmd, err := script.NewCaptureCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Actions = append(scr.Actions, cmd)
		//case script.CmdCopy:
		//	cmd, err := script.NewCopyCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Actions = append(scr.Actions, cmd)
		//case script.CmdRun:
		//	cmd, err := script.NewRunCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Actions = append(scr.Actions, cmd)
		//case script.CmdKubeGet:
		//	cmd, err := script.NewKubeGetCommand(line, rawArgs)
		//	if err != nil {
		//		return nil, err
		//	}
		//	scr.Actions = append(scr.Actions, cmd)
		default:
			return nil, fmt.Errorf("%s not supported", cmdName)
		}
		logrus.Debugf("%s parsed OK", cmdName)
		line++
	}
	logrus.Debug("Done parsing")
	// return enforceDefaults(&scr)
	return scr, nil
}

// enforceDefaults adds missing defaults to the script
//func enforceDefaults(scr *script.Script) (*script.Script, error) {
//	logrus.Debug("Applying default values")
//	if _, ok := scr.Preambles[script.CmdAs]; !ok {
//		cmd, err := script.NewAsCommand(0, fmt.Sprintf("userid:%d groupid:%d", os.Getuid(), os.Getgid()))
//		if err != nil {
//			return scr, err
//		}
//		logrus.Debugf("AS %s:%s (as default)", cmd.GetUserId(), cmd.GetGroupId())
//		scr.Preambles[script.CmdAs] = []script.Directive{cmd}
//	}
//
//	if _, ok := scr.Preambles[script.CmdFrom]; !ok {
//		cmd, err := script.NewFromCommand(0, script.Defaults.FromValue)
//		if err != nil {
//			return nil, err
//		}
//		logrus.Debugf("FROM %v (as default)", cmd.Nodes())
//		scr.Preambles[script.CmdFrom] = []script.Directive{cmd}
//	}
//	if _, ok := scr.Preambles[script.CmdAuthConfig]; !ok {
//		cmd, err := script.NewAuthConfigCommand(0, fmt.Sprintf("username:${USER} private-key:${HOME}/.ssh/id_rsa"))
//		if err != nil {
//			return nil, err
//		}
//		logrus.Debug("AUTHCONFIG set with default")
//		scr.Preambles[script.CmdAuthConfig] = []script.Directive{cmd}
//	}
//	if _, ok := scr.Preambles[script.CmdWorkDir]; !ok {
//		cmd, err := script.NewWorkdirCommand(0, fmt.Sprintf("path:%s", script.Defaults.WorkdirValue))
//		if err != nil {
//			return nil, err
//		}
//		logrus.Debugf("WORKDIR %s (as default)", cmd.Path())
//		scr.Preambles[script.CmdWorkDir] = []script.Directive{cmd}
//	}
//
//	if _, ok := scr.Preambles[script.CmdOutput]; !ok {
//		cmd, err := script.NewOutputCommand(0, fmt.Sprintf("path:%s", script.Defaults.OutputValue))
//		if err != nil {
//			return nil, err
//		}
//		logrus.Debugf("OUTPUT %s (as default)", cmd.Path())
//		scr.Preambles[script.CmdOutput] = []script.Directive{cmd}
//	}
//
//	if _, ok := scr.Preambles[script.CmdKubeConfig]; !ok {
//		cmd, err := script.NewKubeConfigCommand(0, fmt.Sprintf("path:%s", script.Defaults.KubeConfigValue))
//		if err != nil {
//			return nil, err
//		}
//		logrus.Debugf("KUBECONFIG %s (as default)", cmd.Path())
//		scr.Preambles[script.CmdKubeConfig] = []script.Directive{cmd}
//	}
//	return scr, nil
//}

func splitDirectiveLine(line string) (directive, rawArgs string, err error) {
	if len(line) == 0 {
		return "", "", nil
	}
	parts := directiveRegx.FindStringSubmatch(line)
	// len(parts) should be 3
	// [orig string, directiveName, rawArgs]
	if len(parts) != 3 {
		return "", "", fmt.Errorf("malformed directive: %s", line)
	}
	return parts[1], parts[2], nil
}

func isNamedParam(str string) bool {
	return namedParamRegx.MatchString(str)
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
	params, err := commandSplit(rawArgs)
	if err != nil {
		return nil, err
	}

	// for each, split pram:<pram-value> into {param, <param-val>}
	for _, param := range params {
		cmdName, cmdStr, err := namedParamSplit(param)
		if err != nil {
			return nil, fmt.Errorf("map args: %s", err)
		}
		argMap[cmdName] = cmdStr
	}

	return argMap, nil
}

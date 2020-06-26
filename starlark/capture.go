// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/ssh"
)

// captureFunc is a built-in starlark function that runs a provided command and
// captures the result of the command in a specified file stored in workdir.
// If resources and workdir are not provided, captureFunc uses defaults from starlark thread generated
// by previous calls to resources() and crashd_config().
// Starlark format: capture(cmd="command" [,resources=resources][,workdir=path][,file_name=name][,desc=description])
func captureFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var cmdStr string
	if args != nil && args.Len() == 1 {
		cmd, ok := args.Index(0).(starlark.String)
		if !ok {
			return starlark.None, fmt.Errorf("%s: default argument must be a string", identifiers.capture)
		}
		cmdStr = string(cmd)
	}

	// grab named arguments
	var dictionary starlark.StringDict
	if kwargs != nil {
		dict, err := kwargsToStringDict(kwargs)
		if err != nil {
			return starlark.None, err
		}
		dictionary = dict
	}

	if dictionary["cmd"] != nil {
		if cmd, ok := dictionary["cmd"].(starlark.String); ok {
			cmdStr = string(cmd)
		}
	}

	if len(cmdStr) == 0 {
		return starlark.None, fmt.Errorf("%s: missing command string", identifiers.capture)
	}

	var fileName string
	if dictionary["file_name"] != nil {
		if cmd, ok := dictionary["file_name"].(starlark.String); ok {
			fileName = string(cmd)
		}
	}

	var desc string
	if dictionary["desc"] != nil {
		if cmd, ok := dictionary["desc"].(starlark.String); ok {
			desc = string(cmd)
		}
	}

	// extract workdir
	var workdir string
	if dictionary["workdir"] != nil {
		if dir, ok := dictionary["workdir"].(starlark.String); ok {
			workdir = string(dir)
		}
	}
	if len(workdir) == 0 {
		if dir, err := getWorkdirFromThread(thread); err == nil {
			workdir = dir
		}
	}
	if len(workdir) == 0 {
		workdir = defaults.workdir
	}

	// extract resources
	var resources *starlark.List
	if dictionary[identifiers.resources] != nil {
		res, ok := dictionary[identifiers.resources].(*starlark.List)
		if !ok {
			return starlark.None, fmt.Errorf("%s: unexpected resources type", identifiers.capture)
		}
		resources = res
	}
	if resources == nil {
		res := thread.Local(identifiers.resources)
		if res == nil {
			return starlark.None, fmt.Errorf("%s: default resources not found", identifiers.capture)
		}
		resList, ok := res.(*starlark.List)
		if !ok {
			return starlark.None, fmt.Errorf("%s: unexpected resources type", identifiers.capture)
		}
		resources = resList
	}

	results, err := execCapture(cmdStr, workdir, fileName, desc, resources)
	if err != nil {
		return starlark.None, err
	}

	// build list of struct as result
	var resultList []starlark.Value
	for _, result := range results {
		if len(results) == 1 {
			return result.toStarlarkStruct(), nil
		}
		resultList = append(resultList, result.toStarlarkStruct())
	}

	return starlark.NewList(resultList), nil
}

func execCapture(cmdStr, rootPath, fileName, desc string, resources *starlark.List) ([]runResult, error) {
	if resources == nil {
		return nil, fmt.Errorf("%s: missing resources", identifiers.capture)
	}

	logrus.Debugf("%s: executing command on %d resources", identifiers.capture, resources.Len())
	var results []runResult
	for i := 0; i < resources.Len(); i++ {
		val := resources.Index(i)
		res, ok := val.(*starlarkstruct.Struct)
		if !ok {
			return nil, fmt.Errorf("%s: unexpected resource type", identifiers.run)
		}

		val, err := res.Attr("kind")
		if err != nil {
			return nil, fmt.Errorf("%s: resource.kind: %s", identifiers.capture, err)
		}
		kind := val.(starlark.String)

		val, err = res.Attr("transport")
		if err != nil {
			return nil, fmt.Errorf("%s: resource.transport: %s", identifiers.capture, err)
		}
		transport := val.(starlark.String)

		val, err = res.Attr("host")
		if err != nil {
			return nil, fmt.Errorf("%s: resource.host: %s", identifiers.capture, err)
		}
		host := string(val.(starlark.String))
		rootDir := filepath.Join(rootPath, sanitizeStr(host))

		switch {
		case string(kind) == identifiers.hostResource && string(transport) == "ssh":
			result, err := execCaptureSSH(host, cmdStr, rootDir, fileName, desc, res)
			if err != nil {
				logrus.Error(err)
				continue
			}
			results = append(results, result)
		default:
			logrus.Errorf("%s: unsupported or invalid resource kind: %s", identifiers.capture, kind)
			continue
		}
	}

	return results, nil
}

func execCaptureSSH(host, cmdStr, rootDir, fileName, desc string, res *starlarkstruct.Struct) (runResult, error) {
	sshCfg := starlarkstruct.FromKeywords(starlarkstruct.Default, makeDefaultSSHConfig())
	if val, err := res.Attr(identifiers.sshCfg); err == nil {
		if cfg, ok := val.(*starlarkstruct.Struct); ok {
			sshCfg = cfg
		}
	}

	args, err := getSSHArgsFromCfg(sshCfg)
	if err != nil {
		return runResult{}, err
	}
	args.Host = host

	// create dir for the host
	if err := os.MkdirAll(rootDir, 0744); err != nil && !os.IsExist(err) {
		return runResult{}, err
	}

	if len(fileName) == 0 {
		fileName = fmt.Sprintf("%s.txt", sanitizeStr(cmdStr))
	}
	filePath := filepath.Join(rootDir, fileName)

	logrus.Debugf("%s: capturing command on %s using ssh: [%s]", identifiers.capture, args.Host, cmdStr)

	reader, err := ssh.RunRead(args, cmdStr)
	if err != nil {
		if err := captureOutput(strings.NewReader(err.Error()), filePath, fmt.Sprintf("%s: failed", cmdStr)); err != nil {
			return runResult{}, err
		}
	}

	if err := captureOutput(reader, filePath, desc); err != nil {
		return runResult{}, err
	}

	return runResult{resource: args.Host, result: filePath, err: err}, nil

}

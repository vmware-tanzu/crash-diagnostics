// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/ssh"
)

// copyFromFunc is a built-in starlark function that copies file resources from
// specified compute resources and saves them on the local machine
// in subdirectory under workdir.
//
// If resources and workdir are not provided, copyFromFunc uses defaults from starlark thread generated
// by previous calls to resources(), ssh_config, and crashd_config().
//
// Starlark format: copy_from([<path>] [,path=<list>, resources=resources, workdir=path])
func copyFromFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var sourcePath string
	if args != nil && args.Len() == 1 {
		if path, ok := args.Index(0).(starlark.String); ok {
			sourcePath = string(path)
		}
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

	if dictionary["path"] != nil {
		if path, ok := dictionary["path"].(starlark.String); ok {
			sourcePath = string(path)
		}
	}

	if sourcePath == "" {
		return starlark.None, fmt.Errorf("%s: path arg not set", identifiers.copyFrom)
	}

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
		return starlark.None, fmt.Errorf("%s: workdir arg not set", identifiers.copyFrom)
	}

	// extract resources
	var resources *starlark.List
	if dictionary[identifiers.resources] != nil {
		if res, ok := dictionary[identifiers.resources].(*starlark.List); ok {
			resources = res
		}
	}
	if resources == nil {
		res, err := getResourcesFromThread(thread)
		if err != nil {
			return starlark.None, fmt.Errorf("%s: %s", identifiers.copyFrom, err)
		}
		resources = res
	}

	results, err := execCopy(workdir, sourcePath, resources)
	if err != nil {
		return starlark.None, fmt.Errorf("%s: %s", identifiers.copyFrom, err)
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

func execCopy(rootPath string, path string, resources *starlark.List) ([]commandResult, error) {
	if resources == nil {
		return nil, fmt.Errorf("%s: missing resources", identifiers.copyFrom)
	}

	var results []commandResult
	for i := 0; i < resources.Len(); i++ {
		val := resources.Index(i)
		res, ok := val.(*starlarkstruct.Struct)
		if !ok {
			return nil, fmt.Errorf("%s: unexpected resource type", identifiers.copyFrom)
		}

		val, err := res.Attr("kind")
		if err != nil {
			return nil, fmt.Errorf("%s: resource.kind: %s", identifiers.copyFrom, err)
		}
		kind := val.(starlark.String)

		val, err = res.Attr("transport")
		if err != nil {
			return nil, fmt.Errorf("%s: resource.transport: %s", identifiers.copyFrom, err)
		}
		transport := val.(starlark.String)

		val, err = res.Attr("host")
		if err != nil {
			return nil, fmt.Errorf("%s: resource.host: %s", identifiers.copyFrom, err)
		}
		host := string(val.(starlark.String))
		rootDir := filepath.Join(rootPath, sanitizeStr(host))

		switch {
		case string(kind) == identifiers.hostResource && string(transport) == "ssh":
			result, err := execCopySCP(host, rootDir, path, res)
			if err != nil {
				logrus.Errorf("%s: failed to copyFrom %s: %s", identifiers.copyFrom, path, err)
			}
			results = append(results, result)
		default:
			logrus.Errorf("%s: unsupported or invalid resource kind: %s", identifiers.copyFrom, kind)
			continue
		}
	}

	return results, nil
}

func execCopySCP(host, rootDir, path string, res *starlarkstruct.Struct) (commandResult, error) {
	sshCfg := starlarkstruct.FromKeywords(starlarkstruct.Default, makeDefaultSSHConfig())
	if val, err := res.Attr(identifiers.sshCfg); err == nil {
		if cfg, ok := val.(*starlarkstruct.Struct); ok {
			sshCfg = cfg
		}
	}

	args, err := getSSHArgsFromCfg(sshCfg)
	if err != nil {
		return commandResult{}, err
	}
	args.Host = host

	// create dir for the host
	if err := os.MkdirAll(rootDir, 0744); err != nil && !os.IsExist(err) {
		return commandResult{}, err
	}

	err = ssh.CopyFrom(args, rootDir, path)
	return commandResult{resource: args.Host, result: filepath.Join(rootDir, path), err: err}, err
}

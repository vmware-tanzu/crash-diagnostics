// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"os"
	"strings"
)

var (
	kubegetParams = struct {
		containers,
		namespsaces,
		groups,
		kinds,
		versions,
		names,
		labels,
		what string
	}{
		containers:  "containers",
		namespsaces: "namespaces",
		groups:      "groups",
		kinds:       "kinds",
		versions:    "versions",
		names:       "names",
		labels:      "labels",
		what:        "what",
	}
)

// KubeGetCommand represents a KUBEGET directive which can have the following forms:
//     KUBEGET objects namespaces:<namespace> groups:<api-grou-name> kinds:<object-kind> versions:<object-version> names:<object-name> labels:<object-labels>
//     KUBEGET logs namespaces:<namespace> labels:<object-labels> containers:<list-of-containers-names>
//     KUBEGET all labels:<label-list>
// The first param (what-param) is required param that indicates what resource to get
// with valid values of `objects`, `logs` and `all`.  That param can also appear with the name `what`:
//     KUBEGET what:"all" labels:<label-list>
type KubeGetCommand struct {
	cmd
}

// NewKubeGetCommand creates a value of type *KubeGetCommand from a script
func NewKubeGetCommand(index int, rawArgs string) (*KubeGetCommand, error) {
	if err := validateRawArgs(CmdKubeGet, rawArgs); err != nil {
		return nil, err
	}

	// parse the `what` param
	// from rawArgs: <what-param> ... <other-named-params>
	// 1) handle what-param when it's not named:
	var what string
	if !strings.Contains(rawArgs, "what:") {
		// assume first word is what-param
		params := spaceSep.Split(rawArgs, 2)
		if len(params) == 2 {
			rawArgs = params[1]
		} else {
			rawArgs = ""
		}

		switch params[0] {
		case "objects", "logs", "all":
			what = params[0]
		default:
			what = "objects"
		}

		// append named what-param to rawArgs string
		rawArgs = fmt.Sprintf("what:%s %s", what, rawArgs)
	}

	// map remaining params
	var argMap map[string]string
	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("KUBEGET: %v", err)
	}

	cmd := &KubeGetCommand{cmd: cmd{index: index, name: CmdKubeGet, args: argMap}}
	if err := validateCmdArgs(CmdKubeGet, argMap); err != nil {
		return nil, err
	}

	return cmd, nil
}

// Index is the position of the command in the script
func (c *KubeGetCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *KubeGetCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *KubeGetCommand) Args() map[string]string {
	return c.cmd.args
}

// What returns the type of resource to get (i.e. objects, logs, all)
func (c *KubeGetCommand) What() string {
	return os.ExpandEnv(c.args[kubegetParams.what])
}

// Containers returns a comma-sep list of containers from which to retrieve logs
func (c *KubeGetCommand) Containers() string {
	return os.ExpandEnv(c.args[kubegetParams.namespsaces])
}

// Namespaces returns a comma-sep list of namespaces from which to retrieve objects
func (c *KubeGetCommand) Namespaces() string {
	return os.ExpandEnv(c.args[kubegetParams.namespsaces])
}

// Groups returns a comma-sep resource groups from which to retrieve objects
func (c *KubeGetCommand) Groups() string {
	return os.ExpandEnv(c.args[kubegetParams.groups])
}

// Versions returns a comma-sep of resource versions to retrieve
func (c *KubeGetCommand) Versions() string {
	return os.ExpandEnv(c.args[kubegetParams.groups])
}

// Kinds returns a comma-sep list of resource object kinds to retrieve
func (c *KubeGetCommand) Kinds() string {
	return os.ExpandEnv(c.args[kubegetParams.kinds])
}

// Names returns a comma-sep list of resource names to retrieve
func (c *KubeGetCommand) Names() string {
	return os.ExpandEnv(c.args[kubegetParams.names])
}

// Labels returns a comma-sep list of resource labels to match
func (c *KubeGetCommand) Labels() string {
	return os.ExpandEnv(c.args[kubegetParams.kinds])
}

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

// AsCommand represents AS directive:
//
// AS userid:"userid" [groupid:"groupid"]
//
// Param userid required; groupid optional
type AsCommand struct {
	cmd
	user *user.User
}

// NewAsCommand returns *AsCommand with parsed arguments
func NewAsCommand(index int, rawArgs string) (*AsCommand, error) {
	if err := validateRawArgs(CmdAs, rawArgs); err != nil {
		return nil, err
	}

	argMap, err := mapArgs(rawArgs)
	if err != nil {
		return nil, fmt.Errorf("AS: %v", err)
	}
	if err := validateCmdArgs(CmdAs, argMap); err != nil {
		return nil, err
	}
	//validate required param
	if _, ok := argMap["userid"]; len(argMap) == 1 && !ok {
		return nil, fmt.Errorf("AS requires parameter userid")
	}
	// fill in default
	if len(argMap) == 1 {
		argMap["groupid"] = fmt.Sprintf("%d", os.Getgid())
	}

	cmd := &AsCommand{cmd: cmd{index: index, name: CmdAs, args: argMap}}

	return cmd, nil
}

// Index is the position of the command in the script
func (c *AsCommand) Index() int {
	return c.cmd.index
}

// Name represents the name of the command
func (c *AsCommand) Name() string {
	return c.cmd.name
}

// Args returns a slice of raw command arguments
func (c *AsCommand) Args() map[string]string {
	return c.cmd.args
}

// GetUserId returns the userid specified in AS
func (c *AsCommand) GetUserId() string {
	return os.ExpandEnv(c.cmd.args["userid"])
}

// GetGroupId returns the gid specified in AS
func (c *AsCommand) GetGroupId() string {
	return os.ExpandEnv(c.cmd.args["groupid"])
}

// GetCredentials returns the uid and gid value extracted from Args
func (c *AsCommand) GetCredentials() (uid, gid int, err error) {
	if c.user != nil {
		return getUserIds(c.user)
	}

	var usr *user.User
	_, err = strconv.Atoi(c.GetUserId())
	if err != nil { // is id really a username
		usr, err = user.Lookup(c.GetUserId())
		if err != nil {
			return -1, -1, err
		}
	} else {
		usr, err = user.LookupId(c.GetUserId())
		if err != nil {
			return -1, -1, err
		}
	}

	c.user = usr
	return getUserIds(c.user)
}

func getUserIds(usr *user.User) (uid int, gid int, err error) {
	if usr == nil {
		return 0, 0, fmt.Errorf("unable to lookup credentials, user nil")
	}
	uid, err = strconv.Atoi(usr.Uid)
	if err != nil {
		return -1, -1, fmt.Errorf("bad user id %s", usr.Uid)
	}
	gid, err = strconv.Atoi(usr.Gid)
	if err != nil {
		return -1, -1, fmt.Errorf("bad group id %s", usr.Gid)
	}
	return
}

package script

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

type AsCommand struct {
	cmd
	user    *user.User
	userid  string
	groupid string
}

func NewAsCommand(index int, args []string) (*AsCommand, error) {
	cmd := &AsCommand{cmd: cmd{index: index, name: CmdAs, args: args}}

	if err := validateCmdArgs(CmdAs, args); err != nil {
		return nil, err
	}

	if len(args) > 0 {
		asParts := strings.Split(args[0], ":")
		if len(asParts) > 1 {
			cmd.groupid = asParts[1]
		}
		cmd.userid = asParts[0]
	} else {
		cmd.userid = fmt.Sprintf("%d", os.Getuid())
		cmd.groupid = fmt.Sprintf("%d", os.Getgid())
	}

	return cmd, nil
}

func (c *AsCommand) Index() int {
	return c.cmd.index
}

func (c *AsCommand) Name() string {
	return c.cmd.name
}

func (c *AsCommand) Args() []string {
	return c.cmd.args
}

func (c *AsCommand) GetCredentials() (uid, gid int, err error) {
	if c.user != nil {
		return getUserIds(c.user)
	}

	var usr *user.User
	_, err = strconv.Atoi(c.userid)
	if err != nil { // is id really a username
		usr, err = user.Lookup(c.userid)
		if err != nil {
			return -1, -1, err
		}
	} else {
		usr, err = user.LookupId(c.userid)
		if err != nil {
			return -1, -1, err
		}
	}

	c.user = usr
	return getUserIds(c.user)
}

func getUserIds(usr *user.User) (uid int, gid int, err error) {
	if usr == nil {
		return 0, 0, fmt.Errorf("Unable to lookup credentials, user nil")
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

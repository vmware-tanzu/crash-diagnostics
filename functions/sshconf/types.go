// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package sshconf represents the `ssh_config` starlark function
package sshconf

import (
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
)

var (
	DefaultUsername    = func() string { return functions.DefaultUsername() }
	DefaultPort        = func() string { return "22" }
	DefaultMaxRetries  = func() int64 { return 3 }
	DefaultConnTimeout = func() int64 { return 30 }
	DefaultPKPath      = func() string {
		return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}
	DefaultConfig = makeDefaultSSHConfig

	Identifier      = string(Name)
	AgentIdentifier = "crashd_ssh_agent"
)

// Args represent input arguments passed to starlark function.
// Args can also be used as output arguments to built-in function.
//
// The argument map follows:
//   - username - username
//   - port string - SSH port
//   - private_key_path string - SSH private key path
//   - jump_user string - jump host username
//   - jump_host string - jump host name
//   - max_retries [] string - maximum retires for SSH
//   - conn_timeout bool - timeout for connection
//
type Args struct {
	Username       string `name:"username"`
	Port           string `name:"port" optional:"true"`
	PrivateKeyPath string `name:"private_key_path" optional:"true"`
	JumpUsername   string `name:"jump_user" optional:"true"`
	JumpHost       string `name:"jump_host" optional:"true"`
	MaxRetries     int64  `name:"max_retries" optional:"true"`
	ConnTimeout    int64  `name:"conn_timeout" optional:"true"`
}

// Config is a configuration returned by the command function
type Config struct {
	Username       string `name:"username"`
	Port           string `name:"port"`
	PrivateKeyPath string `name:"private_key_path"`
	JumpUsername   string `name:"jump_user"`
	JumpHost       string `name:"jump_host"`
	MaxRetries     int64  `name:"max_retries"`
	ConnTimeout    int64  `name:"conn_timeout"`
}

type Result struct {
	Error  string `name:"error"`
	Config Config `name:"config"`
}

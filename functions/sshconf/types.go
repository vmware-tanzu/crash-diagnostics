// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package sshconf represents the `ssh_config` starlark function
package sshconf

import (
	"os"
	"path/filepath"
)

var (
	DefaultPort        = func() string { return "22" }
	DefaultMaxRetries  = func() int64 { return 3 }
	DefaultConnTimeout = func() int64 { return 30 }
	DefaultPKPath      = func() string {
		return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}

	SSHAgentIdentifier = "crashd_ssh_agent"
)

type Args struct {
	Username       string `name:"username"`
	Port           string `name:"port" optional:"true"`
	PrivateKeyPath string `name:"private_key_path" optional:"true"`
	JumpUsername   string `name:"jump_user" optional:"true"`
	JumpHost       string `name:"jump_host" optional:"true"`
	MaxRetries     int64  `name:"max_retries" optional:"true"`
	ConnTimeout    int64  `name:"conn_timeout" optional:"true"`
}

type Config struct {
	Error          string `name:"error"`
	Username       string `name:"username"`
	Port           string `name:"port"`
	PrivateKeyPath string `name:"private_key_path"`
	JumpUsername   string `name:"jump_user"`
	JumpHost       string `name:"jump_host"`
	MaxRetries     int64  `name:"max_retries"`
	ConnTimeout    int64  `name:"conn_timeout"`
}

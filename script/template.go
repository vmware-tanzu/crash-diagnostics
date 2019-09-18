// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"text/template"

	"github.com/sirupsen/logrus"
)

type templateVars struct {
	Home     string
	Username string
	Pwd      string
}

type templateConfig struct {
	funcs template.FuncMap
	vars  templateVars
}

var (
	store   = make(map[string]string)
	tempCfg = templateConfig{
		funcs: template.FuncMap{
			"get": func(key string) string {
				return store[key]
			},

			"set": func(key, value string) string {
				store[key] = value
				return key
			},
		},

		vars: templateVars{
			Home:     homedir(),
			Username: username(),
			Pwd:      pwd(),
		},
	}
)

func applyTemplate(dest io.Writer, src string) error {
	tmpl, err := template.New("default").Funcs(tempCfg.funcs).Parse(src)
	if err != nil {
		return fmt.Errorf("Failed to apply template: %s", err)
	}
	if err := tmpl.Execute(dest, tempCfg.vars); err != nil {
		return fmt.Errorf("Template execution failed: %s", err)
	}
	return nil
}

func homedir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		logrus.Warnf("Failed to determine user home dir: %s", err)
	}
	return dir
}

func username() string {
	usr, err := user.Current()
	if err != nil {
		logrus.Warnf("Failed to determine user: %s", err)
	}
	return usr.Username
}

func pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Warnf("Failed to determine working directory (pwd) for binary : %s", err)
	}
	return dir
}

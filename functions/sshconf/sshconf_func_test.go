// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sshconf

import (
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
)

func TestSSHConfigFunc(t *testing.T) {
	tests := []struct {
		name   string
		kwargs []starlark.Tuple
		eval   func(*testing.T, []starlark.Tuple)
	}{
		{
			name: "no args",
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err == nil {
					t.Fatal("expecting argument error, got none")
				}
			},
		},
		{
			name:   "with just username",
			kwargs: []starlark.Tuple{{starlark.String("username"), starlark.String("foo")}},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var conf Config
				if err := typekit.Starlark(val).Go(&conf); err != nil {
					t.Fatal(err)
				}
				if conf.Username != "foo" {
					t.Errorf("unexpected username value: %s", conf.Username)
				}
				if conf.Port != "22" {
					t.Errorf("unexpected port value: %s", conf.Port)
				}
				if conf.PrivateKeyPath != DefaultPKPath() {
					t.Errorf("unexpected pk path value: %s", conf.PrivateKeyPath)
				}
			},
		},
		{
			name: "with configs",
			kwargs: []starlark.Tuple{
				{starlark.String("username"), starlark.String("foo")},
				{starlark.String("port"), starlark.String("44")},
				{starlark.String("private_key_path"), starlark.String("./ssh/path")},
				{starlark.String("max_retries"), starlark.MakeInt(32)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var conf Config
				if err := typekit.Starlark(val).Go(&conf); err != nil {
					t.Fatal(err)
				}
				if conf.Username != "foo" {
					t.Errorf("unexpected username value: %s", conf.Username)
				}
				if conf.Port != "44" {
					t.Errorf("unexpected port value: %s", conf.Port)
				}
				if conf.PrivateKeyPath != "./ssh/path" {
					t.Errorf("unexpected pk path value: %s", conf.PrivateKeyPath)
				}
				if conf.MaxRetries != 32 {
					t.Errorf("unexpected max retries value: %d", conf.MaxRetries)
				}
			},
		},
		//{
		//	name: "with shell",
		//	kwargs: []starlark.Tuple{
		//		{starlark.String("default_shell"), starlark.String("/a/b/c")},
		//	},
		//	eval: func(t *testing.T, kwargs []starlark.Tuple) {
		//		val, err := Func(&starlark.Thread{}, nil, nil, kwargs)
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		var conf Result
		//		if err := typekit.Starlark(val).Go(&conf); err != nil {
		//			t.Fatal(err)
		//		}
		//		if conf.DefaultShell != "/a/b/c" {
		//			t.Errorf("unexpected defaultShell value: %s", conf.DefaultShell)
		//		}
		//		if err := os.RemoveAll(conf.Workdir); err != nil {
		//			t.Error(err)
		//		}
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.kwargs)
		})
	}
}

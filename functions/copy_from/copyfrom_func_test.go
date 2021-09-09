// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package copy_from

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers/hostlist"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf/make_sshconf"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestRunFunc(t *testing.T) {
	tests := []struct {
		name        string
		remoteFiles map[string]string
		kwargs      func(*testing.T) []starlark.Tuple
		eval        func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name: "no args",
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				t.Logf("expected err: %s", err)
				if err == nil {
					t.Fatal("expected error, but got nil")
				}
			},
		},
		{
			name: "missing resources",
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshConf := sshconf.DefaultConfig()
				sshArg, err := functions.Result(make_sshconf.Name, sshConf)
				if err != nil {
					t.Fatal(err)
				}
				return []starlark.Tuple{
					{starlark.String("path"), starlark.String("/path/to/copy")},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				t.Logf("expected err: %s", err)
				if err == nil {
					t.Fatal("expecting error, but got nil")
				}
			},
		},
		{
			name: "missing sshconf",
			kwargs: func(t *testing.T) []starlark.Tuple {
				resources := providers.Resources{
					Provider: string(hostlist.Name),
					Hosts:    []string{"127.0.0.1"},
				}
				resArg, err := functions.AsStarlarkStruct(resources)
				if err != nil {
					t.Fatal(err)
				}
				return []starlark.Tuple{
					{starlark.String("path"), starlark.String("/path/to/copy")},
					{starlark.String("resources"), resArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				res, err := Func(&starlark.Thread{}, nil, nil, kwargs)
				if err != nil {
					t.Fatal("unexpected error:", err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatal("unexpected function error: ", result.Error)
				}
			},
		},
		{
			name:        "single machine single file",
			remoteFiles: map[string]string{"foo.txt": "FooBar"},
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"username":         starlark.String(testSupport.CurrentUsername()),
					"port":             starlark.String(testSupport.PortValue()),
					"private_key_path": starlark.String(testSupport.PrivateKeyPath()),
					"max_retries":      starlark.MakeInt(testSupport.MaxConnectionRetries()),
				})
				resArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"provider": starlark.String(hostlist.Name),
					"hosts":    starlark.NewList([]starlark.Value{starlark.String("127.0.0.1")}),
				})
				return []starlark.Tuple{
					{starlark.String("path"), starlark.String("foo.txt")},
					{starlark.String("resources"), resArg},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				thread := &starlark.Thread{}
				res, err := Func(thread, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}
				for _, copy := range result.Copies {
					if copy.Error != "" {
						t.Errorf("copy_from failed: host %s: %s", copy.Host, copy.Error)
					}
					if _, err := os.Stat(copy.Path); err != nil {
						t.Errorf("copy_from file not found: %s", err)
					}
					if err := os.RemoveAll(copy.Path); err != nil {
						t.Errorf("copy_from_test: failed to remove file: %s", err)
					}
				}
			},
		},
		{
			name:        "multiple machines multiple source files",
			remoteFiles: map[string]string{"bar.txt": "BarBar", "foo.txt": "FooBar", "baz.txt": "BazBuz"},
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"username":         starlark.String(testSupport.CurrentUsername()),
					"port":             starlark.String(testSupport.PortValue()),
					"private_key_path": starlark.String(testSupport.PrivateKeyPath()),
					"max_retries":      starlark.MakeInt(testSupport.MaxConnectionRetries()),
				})
				resArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"provider": starlark.String(hostlist.Name),
					"hosts":    starlark.NewList([]starlark.Value{starlark.String("localhost"), starlark.String("127.0.0.1")}),
				})
				return []starlark.Tuple{
					{starlark.String("path"), starlark.String("baz.txt")},
					{starlark.String("resources"), resArg},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				thread := &starlark.Thread{}
				res, err := Func(thread, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}
				for _, copy := range result.Copies {
					if copy.Error != "" {
						t.Errorf("copy_from failed: host %s: %s", copy.Host, copy.Error)
					}
					if _, err := os.Stat(copy.Path); err != nil {
						t.Errorf("copy_from file not found: %s", err)
					}
					if err := os.RemoveAll(copy.Path); err != nil {
						t.Errorf("copy_from_test: failed to remove file: %s", err)
					}
				}
			},
		},
		{
			name:        "single machine nested source files",
			remoteFiles: map[string]string{"bar/bar.txt": "BarBar", "bar/foo.txt": "FooBar"},
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"username":         starlark.String(testSupport.CurrentUsername()),
					"port":             starlark.String(testSupport.PortValue()),
					"private_key_path": starlark.String(testSupport.PrivateKeyPath()),
					"max_retries":      starlark.MakeInt(testSupport.MaxConnectionRetries()),
				})
				resArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"provider": starlark.String(hostlist.Name),
					"hosts":    starlark.NewList([]starlark.Value{starlark.String("127.0.0.1")}),
				})
				return []starlark.Tuple{
					{starlark.String("path"), starlark.String("bar/foo.txt")},
					{starlark.String("resources"), resArg},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				thread := &starlark.Thread{}
				res, err := Func(thread, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}
				for _, copy := range result.Copies {
					if copy.Error != "" {
						t.Errorf("copy_from failed: host %s: %s", copy.Host, copy.Error)
					}
					if _, err := os.Stat(copy.Path); err != nil {
						t.Errorf("copy_from file not found: %s", err)
					}
					if err := os.RemoveAll(copy.Path); err != nil {
						t.Errorf("copy_from_test: failed to remove file: %s", err)
					}
				}
			},
		},
		{
			name:        "single machine gob source files",
			remoteFiles: map[string]string{"bar/bar.txt": "BarBar", "bar/foo.txt": "FooBar"},
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"username":         starlark.String(testSupport.CurrentUsername()),
					"port":             starlark.String(testSupport.PortValue()),
					"private_key_path": starlark.String(testSupport.PrivateKeyPath()),
					"max_retries":      starlark.MakeInt(testSupport.MaxConnectionRetries()),
				})
				resArg := starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
					"provider": starlark.String(hostlist.Name),
					"hosts":    starlark.NewList([]starlark.Value{starlark.String("127.0.0.1")}),
				})
				return []starlark.Tuple{
					{starlark.String("path"), starlark.String("bar/*.txt")},
					{starlark.String("resources"), resArg},
					{starlark.String("ssh_config"), sshArg},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				thread := &starlark.Thread{}
				res, err := Func(thread, nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}
				var result Result
				if err := typekit.Starlark(res).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}

				for _, copy := range result.Copies {
					if copy.Error != "" {
						t.Errorf("copy_from failed: host %s: %s", copy.Host, copy.Error)
					}
					copyDir := filepath.Dir(copy.Path)
					infos, err := ioutil.ReadDir(copyDir)
					if err != nil {
						t.Error(err)
					}
					if len(infos) != 2 {
						t.Error("unexpected number of files copied")
					}
					if err := os.RemoveAll(copyDir); err != nil {
						t.Errorf("copy_from_test: failed to remove file: %s", err)
					}
				}
			},
		},
	}

	sshArgs := ssh.SSHArgs{
		User:           testSupport.CurrentUsername(),
		Host:           "127.0.0.1",
		Port:           testSupport.PortValue(),
		PrivateKeyPath: testSupport.PrivateKeyPath(),
		MaxRetries:     100,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var kwargs []starlark.Tuple
			if test.kwargs != nil {
				kwargs = test.kwargs(t)
			}
			if test.remoteFiles != nil {
				for file, content := range test.remoteFiles {
					ssh.MakeRemoteTestSSHFile(t, sshArgs, file, content)
				}
				defer func() {
					for file := range test.remoteFiles {
						ssh.RemoveRemoteTestSSHFile(t, sshArgs, file)
					}
				}()
			}
			test.eval(t, kwargs)
		})
	}
}

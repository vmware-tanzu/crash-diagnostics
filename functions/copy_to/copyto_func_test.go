// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package copy_to

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions/providers/hostlist"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestRunFunc(t *testing.T) {
	sshArgs := ssh.SSHArgs{
		User:           testSupport.CurrentUsername(),
		Host:           "127.0.0.1",
		Port:           testSupport.PortValue(),
		PrivateKeyPath: testSupport.PrivateKeyPath(),
		MaxRetries:     150,
	}

	tests := []struct {
		name       string
		localFiles map[string]string
		kwargs     func(*testing.T) []starlark.Tuple
		eval       func(t *testing.T, kwargs []starlark.Tuple)
	}{
		//{
		//	name: "no args",
		//	eval: func(t *testing.T, kwargs []starlark.Tuple) {
		//		_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
		//		t.Logf("expected err: %s", err)
		//		if err == nil {
		//			t.Fatal("expected error, but got nil")
		//		}
		//	},
		//},
		//{
		//	name: "missing resources",
		//	kwargs: func(t *testing.T) []starlark.Tuple {
		//		sshConf := sshconf.DefaultConfig()
		//		sshArg, err := functions.Result(make_sshconf.Name, sshConf)
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		return []starlark.Tuple{
		//			{starlark.String("source_path"), starlark.String("/path/to/copy")},
		//			{starlark.String("ssh_config"), sshArg},
		//		}
		//	},
		//	eval: func(t *testing.T, kwargs []starlark.Tuple) {
		//		_, err := Func(&starlark.Thread{}, nil, nil, kwargs)
		//		t.Logf("expected err: %s", err)
		//		if err == nil {
		//			t.Fatal("expecting error, but got nil")
		//		}
		//	},
		//},
		//{
		//	name: "missing sshconf",
		//	kwargs: func(t *testing.T) []starlark.Tuple {
		//		resources := providers.Resources{
		//			Provider: string(hostlist.Name),
		//			Hosts:    []string{"127.0.0.1"},
		//		}
		//		resArg, err := functions.AsStarlarkStruct(resources)
		//		if err != nil {
		//			t.Fatal(err)
		//		}
		//		return []starlark.Tuple{
		//			{starlark.String("source_path"), starlark.String("/path/to/copy")},
		//			{starlark.String("resources"), resArg},
		//		}
		//	},
		//	eval: func(t *testing.T, kwargs []starlark.Tuple) {
		//		res, err := Func(&starlark.Thread{}, nil, nil, kwargs)
		//		if err != nil {
		//			t.Fatal("unexpected error:", err)
		//		}
		//		var result Result
		//		if err := typekit.Starlark(res).Go(&result); err != nil {
		//			t.Fatal(err)
		//		}
		//		if result.Error != "" {
		//			t.Fatal("unexpected function error: ", result.Error)
		//		}
		//	},
		//},
		{
			name:       "single machine single file",
			localFiles: map[string]string{"foo.txt": "FooBar"},
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
					{starlark.String("source_path"), starlark.String("foo.txt")},
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
						t.Fatalf("copy_to failed: host %s: %s", copy.Host, copy.Error)
					}
					ssh.AssertRemoteTestSSHFile(t, sshArgs, copy.TargetPath)

					if err := os.RemoveAll(copy.SourcePath); err != nil {
						t.Errorf("copy_to_test: failed to remove file: %s", err)
					}
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var kwargs []starlark.Tuple
			if test.kwargs != nil {
				kwargs = test.kwargs(t)
			}

			if test.localFiles != nil {
				for file, content := range test.localFiles {
					localFile := filepath.Join(testSupport.TmpDirRoot(), file)
					ssh.MakeLocalTestFile(t, localFile, content)
					ssh.MakeRemoteTestSSHDir(t, sshArgs, file) // if needed.
				}

				defer func() {
					for file := range test.localFiles {
						localFile := filepath.Join(testSupport.TmpDirRoot(), file)
						ssh.RemoveRemoteTestSSHFile(t, sshArgs, file)
						ssh.RemoveLocalTestFile(t, localFile)
					}
				}()
			}

			test.eval(t, kwargs)
		})
	}
}

// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/vmware-tanzu/crash-diagnostics/ssh"
)

func testCopyToFuncForHostResources(t *testing.T, port, privateKey, username string) {
	tests := []struct {
		name       string
		localFiles map[string]string
		kwargs     func(t *testing.T) []starlark.Tuple
		eval       func(t *testing.T, kwargs []starlark.Tuple, sshArgs ssh.SSHArgs)
	}{
		{
			name:       "single machine single file",
			localFiles: map[string]string{"foo.txt": "FooBar"},
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(privateKey, port, username)
				resources := starlark.NewList([]starlark.Value{
					makeTestSSHHostResource("127.0.0.1", sshCfg),
				})
				localFile := starlark.String(filepath.Join(testSupport.TmpDirRoot(), "foo.txt"))
				remoteFile := starlark.String("foo.txt")
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("resources"), resources},
					[]starlark.Value{starlark.String("source_path"), localFile},
					[]starlark.Value{starlark.String("target_path"), remoteFile},
				}
			},

			eval: func(t *testing.T, kwargs []starlark.Tuple, sshArgs ssh.SSHArgs) {
				val, err := copyToFunc(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				var cpErr string
				var targetPath string
				if strct, ok := val.(*starlarkstruct.Struct); ok {
					t.Logf(" starlarkstruct [%#v]", strct)

					if val, err := strct.Attr("err"); err == nil {
						if r, ok := val.(starlark.String); ok {
							cpErr = string(r)
						}
					}
					if val, err := strct.Attr("result"); err == nil {
						if r, ok := val.(starlark.String); ok {
							targetPath = string(r)
						}
					}
				}

				ssh.AssertRemoteTestSSHFile(t, sshArgs, targetPath)

				if cpErr != "" {
					t.Fatal(cpErr)
				}
			},
		},

		{
			name:       "multiple machines single files",
			localFiles: map[string]string{"bar/bar.txt": "BarBar", "bar/foo.txt": "FooBar", "baz.txt": "BazBuz"},
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(privateKey, port, username)
				resources := starlark.NewList([]starlark.Value{
					makeTestSSHHostResource("localhost", sshCfg),
					makeTestSSHHostResource("127.0.0.1", sshCfg),
				})

				localFile := starlark.String(filepath.Join(testSupport.TmpDirRoot(), "bar/bar.txt"))
				remoteFile := starlark.String("bar/bar.txt")

				return []starlark.Tuple{
					[]starlark.Value{starlark.String("resources"), resources},
					[]starlark.Value{starlark.String("source_path"), localFile},
					[]starlark.Value{starlark.String("target_path"), remoteFile},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple, sshArgs ssh.SSHArgs) {
				val, err := copyToFunc(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				resultList, ok := val.(*starlark.List)
				if !ok {
					t.Fatalf("expecting type *starlark.List, got %T", val)
				}

				for i := 0; i < resultList.Len(); i++ {

					var cpErr string
					var targetPath string
					if strct, ok := resultList.Index(i).(*starlarkstruct.Struct); ok {
						if val, err := strct.Attr("err"); err == nil {
							if r, ok := val.(starlark.String); ok {
								cpErr = string(r)
							}
						}
						if val, err := strct.Attr("result"); err == nil {
							if r, ok := val.(starlark.String); ok {
								targetPath = string(r)
							}
						}
					}

					ssh.AssertRemoteTestSSHFile(t, sshArgs, targetPath)

					if cpErr != "" {
						t.Fatal(cpErr)
					}
				}
			},
		},

		{
			name:       "multiple machines glob path",
			localFiles: map[string]string{"bar/bar.txt": "BarBar", "bar/foo.txt": "FooBar", "bar/baz.csv": "BizzBuzz"},
			kwargs: func(t *testing.T) []starlark.Tuple {
				sshCfg := makeTestSSHConfig(privateKey, port, username)
				resources := starlark.NewList([]starlark.Value{
					makeTestSSHHostResource("127.0.0.1", sshCfg),
					makeTestSSHHostResource("localhost", sshCfg),
				})
				localFile := starlark.String(filepath.Join(testSupport.TmpDirRoot(), "bar/baz.csv"))
				remoteFile := starlark.String("bar/baz.cvs")

				return []starlark.Tuple{
					[]starlark.Value{starlark.String("resources"), resources},
					[]starlark.Value{starlark.String("source_path"), localFile},
					[]starlark.Value{starlark.String("target_path"), remoteFile},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple, sshArgs ssh.SSHArgs) {
				val, err := copyToFunc(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatal(err)
				}

				resultList, ok := val.(*starlark.List)
				if !ok {
					t.Fatalf("expecting type *starlark.List, got %T", val)
				}

				for i := 0; i < resultList.Len(); i++ {

					var cpErr string
					var targetPath string
					if strct, ok := resultList.Index(i).(*starlarkstruct.Struct); ok {
						if val, err := strct.Attr("err"); err == nil {
							if r, ok := val.(starlark.String); ok {
								cpErr = string(r)
							}
						}
						if val, err := strct.Attr("result"); err == nil {
							if r, ok := val.(starlark.String); ok {
								targetPath = string(r)
							}
						}
					}

					ssh.AssertRemoteTestSSHFile(t, sshArgs, targetPath)

					if cpErr != "" {
						t.Fatal(cpErr)
					}
				}
			},
		},
	}

	sshArgs := ssh.SSHArgs{User: username, Host: "127.0.0.1", Port: port, PrivateKeyPath: privateKey}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for file, content := range test.localFiles {
				localFile := filepath.Join(testSupport.TmpDirRoot(), file)
				ssh.MakeRemoteTestSSHDir(t, sshArgs, file) // if needed.
				ssh.MakeLocalTestFile(t, localFile, content)
			}
			defer func() {
				for file := range test.localFiles {
					localFile := filepath.Join(testSupport.TmpDirRoot(), file)
					ssh.RemoveRemoteTestSSHFile(t, sshArgs, file)
					ssh.RemoveLocalTestFile(t, localFile)
				}
			}()

			test.eval(t, test.kwargs(t), sshArgs)
		})
	}
}

func testCopyToFuncScriptForHostResources(t *testing.T, port, privateKey, username string) {
	tests := []struct {
		name       string
		localFiles map[string]string
		script     string
		eval       func(t *testing.T, sshArgs ssh.SSHArgs, script string)
	}{
		{
			name:       "multiple machines single copyTo",
			localFiles: map[string]string{"foobar.c": "footext", "bar/bar.txt": "BarBar", "bar/foo.txt": "FooBar", "bar/baz.csv": "BizzBuzz"},
			script: fmt.Sprintf(`
set_defaults(resources(provider = host_list_provider(hosts=["127.0.0.1","localhost"], ssh_config = ssh_config(username="%s", port="%s", private_key_path="%s"))))
result = copy_to(source_path="%s/bar/foo.txt", target_path="bar/foo.txt")`,
				username, port, privateKey, testSupport.TmpDirRoot()),
			eval: func(t *testing.T, sshArgs ssh.SSHArgs, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("copy_to() should be assigned to a variable")
				}
				resultList, ok := resultVal.(*starlark.List)
				if !ok {
					t.Fatalf("expecting type *starlark.List, got %T", resultVal)
				}

				for i := 0; i < resultList.Len(); i++ {

					var cpErr string
					var targetPath string
					if strct, ok := resultList.Index(i).(*starlarkstruct.Struct); ok {
						if val, err := strct.Attr("err"); err == nil {
							if r, ok := val.(starlark.String); ok {
								cpErr = string(r)
							}
						}
						if val, err := strct.Attr("result"); err == nil {
							if r, ok := val.(starlark.String); ok {
								targetPath = string(r)
							}
						}
					}

					ssh.AssertRemoteTestSSHFile(t, sshArgs, targetPath)

					if cpErr != "" {
						t.Fatal(cpErr)
					}
				}
			},
		},

		{
			name:       "resource loop",
			localFiles: map[string]string{"bar/bar.txt": "BarBar", "bar/foo.txt": "FooBar", "bar/baz.csv": "BizzBuzz"},
			script: fmt.Sprintf(`
# execute cmd on each host
def cp(hosts):
	result = []
	for host in hosts:
		result.append(copy_to(source_path="%s/bar/foo.txt", target_path="bar", resources=[host]))
	return result

# configuration
set_defaults(ssh_config(username="%s", port="%s", private_key_path="%s"))
hosts = resources(provider=host_list_provider(hosts=["127.0.0.1","localhost"]))
result = cp(hosts)`, testSupport.TmpDirRoot(), username, port, privateKey),
			eval: func(t *testing.T, sshArgs ssh.SSHArgs, script string) {
				exe := New()
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}

				resultVal := exe.result["result"]
				if resultVal == nil {
					t.Fatal("capture() should be assigned to a variable")
				}
				resultList, ok := resultVal.(*starlark.List)
				if !ok {
					t.Fatalf("expecting type *starlark.List, got %T", resultVal)
				}

				for i := 0; i < resultList.Len(); i++ {

					var cpErr string
					var targetPath string
					if strct, ok := resultList.Index(i).(*starlarkstruct.Struct); ok {
						if val, err := strct.Attr("err"); err == nil {
							if r, ok := val.(starlark.String); ok {
								cpErr = string(r)
							}
						}
						if val, err := strct.Attr("result"); err == nil {
							if r, ok := val.(starlark.String); ok {
								targetPath = string(r)
							}
						}
					}

					ssh.AssertRemoteTestSSHFile(t, sshArgs, targetPath)

					if cpErr != "" {
						t.Fatal(cpErr)
					}
				}
			},
		},
	}

	sshArgs := ssh.SSHArgs{User: username, Host: "127.0.0.1", Port: port, PrivateKeyPath: privateKey}
	for _, test := range tests {
		for file, content := range test.localFiles {
			localFile := filepath.Join(testSupport.TmpDirRoot(), file)
			ssh.MakeRemoteTestSSHDir(t, sshArgs, file) // if needed.
			ssh.MakeLocalTestFile(t, localFile, content)
		}

		defer func() {
			for file := range test.localFiles {
				localFile := filepath.Join(testSupport.TmpDirRoot(), file)
				ssh.RemoveRemoteTestSSHFile(t, sshArgs, file)
				ssh.RemoveLocalTestFile(t, localFile)
			}
		}()

		t.Run(test.name, func(t *testing.T) {
			test.eval(t, sshArgs, test.script)
		})
	}
}

func TestCopyToFuncSSHAll(t *testing.T) {
	port := testSupport.PortValue()
	username := testSupport.CurrentUsername()
	privateKey := testSupport.PrivateKeyPath()

	tests := []struct {
		name string
		test func(t *testing.T, port, privateKey, username string)
	}{
		//{name: "copyToFunc for host resources", test: testCopyToFuncForHostResources},
		{name: "copy_from script for host resources", test: testCopyToFuncScriptForHostResources},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.test(t, port, privateKey, username)
			defer os.RemoveAll(defaults.workdir)
		})
	}
}

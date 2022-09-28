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
)

func TestKubeCapture(t *testing.T) {
	tests := []struct {
		name   string
		kwargs func(t *testing.T) []starlark.Tuple
		eval   func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name: "simple test with namespaced objects",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("what"), starlark.String("objects")},
					[]starlark.Value{starlark.String("groups"), starlark.NewList([]starlark.Value{starlark.String("core")})},
					[]starlark.Value{starlark.String("kinds"), starlark.NewList([]starlark.Value{starlark.String("services")})},
					[]starlark.Value{starlark.String("namespaces"), starlark.NewList([]starlark.Value{starlark.String("default"), starlark.String("kube-system")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeCaptureFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				fileVal, err := resultStruct.Attr("file")
				if err != nil {
					t.Error(err)
				}
				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}

				workDir := fileValStr.GoString()
				fileInfo, err := os.Stat(workDir)
				if err != nil {
					t.Fatalf("stat(%s) failed: %s", workDir, err)
				}
				if !fileInfo.IsDir() {
					t.Fatalf("expecting starlark function to return a dir")
				}
				defer os.RemoveAll(workDir)

				path := filepath.Join(workDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}
				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}

				path = filepath.Join(workDir, "core_v1", "default")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to exist: %s", path, err)
				}
				files, err = os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}
			},
		},
		{
			name: "test for non-namespaced objects",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("what"), starlark.String("objects")},
					[]starlark.Value{starlark.String("groups"), starlark.NewList([]starlark.Value{starlark.String("core")})},
					[]starlark.Value{starlark.String("kinds"), starlark.NewList([]starlark.Value{starlark.String("nodes")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeCaptureFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				fileVal, err := resultStruct.Attr("file")
				if err != nil {
					t.Error(err)
				}
				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}

				captureDir := fileValStr.GoString()
				defer os.RemoveAll(captureDir)

				path := filepath.Join(captureDir, "core_v1")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s file to exist: %s", path, err)
				}
				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) != 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}
				if files[0].IsDir() {
					t.Errorf("expecting to find a regular file, but found dir: %s", files[0].Name())
				}
				if !strings.Contains(files[0].Name(), "nodes-") {
					t.Errorf("expecting to find a node output file, but fond: %s", files[0].Name())
				}

			},
		},
		{
			name: "simple test with objects in categories",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("what"), starlark.String("objects")},
					[]starlark.Value{starlark.String("groups"), starlark.NewList([]starlark.Value{starlark.String("core")})},
					[]starlark.Value{starlark.String("categories"), starlark.NewList([]starlark.Value{starlark.String("all")})},
					[]starlark.Value{starlark.String("namespaces"), starlark.NewList([]starlark.Value{starlark.String("default"), starlark.String("kube-system")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeCaptureFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				fileVal, err := resultStruct.Attr("file")
				if err != nil {
					t.Error(err)
				}
				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}

				captureDir := fileValStr.GoString()
				defer os.RemoveAll(captureDir)

				fileInfo, err := os.Stat(captureDir)
				if err != nil {
					t.Fatalf("stat(%s) failed: %s", captureDir, err)
				}
				if !fileInfo.IsDir() {
					t.Fatalf("expecting starlark function to return a dir")
				}
				path := filepath.Join(captureDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}

				path = filepath.Join(captureDir, "core_v1", "default")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to exist: %s", path, err)
				}

				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}
			},
		},
		{
			name: "search for all logs in a namespace",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("what"), starlark.String("logs")},
					[]starlark.Value{starlark.String("namespaces"), starlark.NewList([]starlark.Value{starlark.String("kube-system")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeCaptureFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				fileVal, err := resultStruct.Attr("file")
				if err != nil {
					t.Error(err)
				}
				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}

				captureDir := fileValStr.GoString()
				if _, err := os.Stat(captureDir); err != nil {
					t.Fatalf("stat(%s) failed: %s", captureDir, err)
				}
				defer os.RemoveAll(captureDir)

				path := filepath.Join(captureDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}

				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 3 {
					t.Error("unexpected number of log files for namespace kube-system:", len(files))
				}
			},
		},
		{
			name: "search for logs for a specific container",
			kwargs: func(t *testing.T) []starlark.Tuple {
				return []starlark.Tuple{
					[]starlark.Value{starlark.String("what"), starlark.String("logs")},
					[]starlark.Value{starlark.String("namespaces"), starlark.NewList([]starlark.Value{starlark.String("kube-system")})},
					[]starlark.Value{starlark.String("containers"), starlark.NewList([]starlark.Value{starlark.String("etcd")})},
				}
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val, err := KubeCaptureFn(newTestThreadLocal(t), nil, nil, kwargs)
				if err != nil {
					t.Fatalf("failed to execute: %s", err)
				}
				resultStruct, ok := val.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("expecting type *starlarkstruct.Struct, got %T", val)
				}

				errVal, err := resultStruct.Attr("error")
				if err != nil {
					t.Error(err)
				}
				resultErr := errVal.(starlark.String).GoString()
				if resultErr != "" {
					t.Fatalf("starlark func failed: %s", resultErr)
				}

				fileVal, err := resultStruct.Attr("file")
				if err != nil {
					t.Error(err)
				}
				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}

				captureDir := fileValStr.GoString()
				if _, err := os.Stat(captureDir); err != nil {
					t.Fatalf("stat(%s) failed: %s", captureDir, err)
				}
				defer os.RemoveAll(captureDir)

				path := filepath.Join(captureDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}

				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) == 0 {
					t.Error("unexpected number of log files for namespace kube-system:", len(files))
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.kwargs(t))
		})
	}
}

func TestKubeCaptureScript(t *testing.T) {
	workdir := testSupport.TmpDirRoot()
	k8sconfig := testSupport.KindKubeConfigFile()
	clusterCtxName := testSupport.KindClusterContextName()

	execute := func(t *testing.T, script string) *starlarkstruct.Struct {
		executor := New()
		if err := executor.Exec("test.kube.capture", strings.NewReader(script)); err != nil {
			t.Fatalf("failed to exec: %s", err)
		}
		if !executor.result.Has("kube_data") {
			t.Fatalf("script result must be assigned to a value")
		}

		data, ok := executor.result["kube_data"].(*starlarkstruct.Struct)
		if !ok {
			t.Fatal("script result is not a struct")
		}
		return data
	}

	tests := []struct {
		name   string
		script string
		eval   func(t *testing.T, script string)
	}{
		{
			name: "simple search with namespaced objects with cluster context",
			script: fmt.Sprintf(`
crashd_config(workdir="%s")
set_defaults(kube_config(path="%s", cluster_context="%s"))
kube_data = kube_capture(what="objects", groups=["core"], kinds=["services"], namespaces=["default", "kube-system"])`, workdir, k8sconfig, clusterCtxName),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				fileVal, err := data.Attr("file")
				if err != nil {
					t.Error(err)
				}

				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}

				workDir := fileValStr.GoString()
				fileInfo, err := os.Stat(workDir)
				if err != nil {
					t.Fatalf("stat(%s) failed: %s", workDir, err)
				}
				if !fileInfo.IsDir() {
					t.Fatalf("expecting starlark function to return a dir")
				}
				defer os.RemoveAll(workDir)

				path := filepath.Join(workDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}
				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}

				path = filepath.Join(workDir, "core_v1", "default")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to exist", path)
				}
				files, err = os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}
			},
		},
		{
			name: "search with non-namespaced objects",
			script: fmt.Sprintf(`
crashd_config(workdir="%s")
set_defaults(kube_config(path="%s"))
kube_data = kube_capture(what="objects", groups=["core"], kinds=["nodes"])`, workdir, k8sconfig),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				fileVal, err := data.Attr("file")
				if err != nil {
					t.Error(err)
				}

				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}

				captureDir := fileValStr.GoString()
				defer os.RemoveAll(captureDir)

				path := filepath.Join(captureDir, "core_v1")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s file to exist", path)
				}
				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) != 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}
				if files[0].IsDir() {
					t.Errorf("expecting to find a regular file, but found dir: %s", files[0].Name())
				}
				if !strings.Contains(files[0].Name(), "nodes-") {
					t.Errorf("expecting to find a node output file, but fond: %s", files[0].Name())
				}
			},
		},
		{
			name: "search with objects in categories",
			script: fmt.Sprintf(`
crashd_config(workdir="%s")
set_defaults(kube_config(path="%s"))
kube_data = kube_capture(what="objects", groups=["core"], categories=["all"], namespaces=["default","kube-system"])`, workdir, k8sconfig),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				fileVal, err := data.Attr("file")
				if err != nil {
					t.Error(err)
				}

				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				captureDir := fileValStr.GoString()
				defer os.RemoveAll(captureDir)

				fileInfo, err := os.Stat(captureDir)
				if err != nil {
					t.Fatalf("stat(%s) failed: %s", captureDir, err)
				}
				if !fileInfo.IsDir() {
					t.Fatalf("expecting starlark function to return a dir")
				}
				path := filepath.Join(captureDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}

				path = filepath.Join(captureDir, "core_v1", "default")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to exist", path)
				}

				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 1 {
					t.Errorf("directory should have at least 1 file but has none: %s:", path)
				}
			},
		},
		{
			name: "search for all logs in a namespace with cluster context",
			script: fmt.Sprintf(`
crashd_config(workdir="%s")
set_defaults(kube_config(path="%s", cluster_context="%s"))
kube_data = kube_capture(what="logs", namespaces=["kube-system"])`, workdir, k8sconfig, clusterCtxName),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				fileVal, err := data.Attr("file")
				if err != nil {
					t.Error(err)
				}

				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				captureDir := fileValStr.GoString()
				if _, err := os.Stat(captureDir); err != nil {
					t.Fatalf("stat(%s) failed: %s", captureDir, err)
				}
				defer os.RemoveAll(captureDir)

				path := filepath.Join(captureDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}

				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) < 3 {
					t.Error("unexpected number of log files for namespace kube-system:", len(files))
				}
			},
		},
		{
			name: "search for logs in specific container",
			script: fmt.Sprintf(`
crashd_config(workdir="%s")
set_defaults(kube_config(path="%s"))
kube_data = kube_capture(what="logs", namespaces=["kube-system"], containers=["etcd"])`, workdir, k8sconfig),
			eval: func(t *testing.T, script string) {
				data := execute(t, script)

				fileVal, err := data.Attr("file")
				if err != nil {
					t.Error(err)
				}

				fileValStr, ok := fileVal.(starlark.String)
				if !ok {
					t.Fatalf("unexpected type for starlark value")
				}
				captureDir := fileValStr.GoString()
				if _, err := os.Stat(captureDir); err != nil {
					t.Fatalf("stat(%s) failed: %s", captureDir, err)
				}
				defer os.RemoveAll(captureDir)

				path := filepath.Join(captureDir, "core_v1", "kube-system")
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expecting %s to be a directory", path)
				}

				files, err := os.ReadDir(path)
				if err != nil {
					t.Fatalf("ReadeDir(%s) failed: %s", path, err)
				}
				if len(files) == 0 {
					t.Error("unexpected number of log files for namespace kube-system:", len(files))
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script)
		})
	}
}

func TestKubeCapture_WithBadKubeconfig(t *testing.T) {
	script := `
cfg = kube_config(path="/foo/bar")
kube_capture(what="logs", namespaces=["kube-system"], containers=["etcd"], kube_config=cfg)
`
	executor := New()
	if err := executor.Exec("test.kube.capture", strings.NewReader(script)); err == nil {
		t.Fatalf("expected failure, but did not get it")
	}
}

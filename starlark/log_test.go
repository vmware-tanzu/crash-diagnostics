// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestLogFunc(t *testing.T) {
	tests := []struct {
		name   string
		kwargs []starlark.Tuple
		test   func(*testing.T, []starlark.Tuple)
	}{
		{
			name: "logging with no prefix",
			kwargs: []starlark.Tuple{
				[]starlark.Value{starlark.String("msg"), starlark.String("Logging from starlark")},
			},
			test: func(t *testing.T, kwargs []starlark.Tuple) {
				var buf bytes.Buffer
				logger := log.New(&buf, "", log.Lshortfile)
				thread := newTestThreadLocal(t)
				thread.SetLocal(identifiers.log, logger)

				if _, err := logFunc(thread, nil, nil, kwargs); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(buf.String(), "Logging from starlark") {
					t.Error("logger has unexpected log msg: ", buf.String())
				}
			},
		},
		{
			name: "logging with prefix",
			kwargs: []starlark.Tuple{
				[]starlark.Value{starlark.String("msg"), starlark.String("Logging from starlark with prefix")},
				[]starlark.Value{starlark.String("prefix"), starlark.String("INFO")},
			},
			test: func(t *testing.T, kwargs []starlark.Tuple) {
				var buf bytes.Buffer
				logger := log.New(&buf, "", log.Lshortfile)
				thread := newTestThreadLocal(t)
				thread.SetLocal(identifiers.log, logger)

				if _, err := logFunc(thread, nil, nil, kwargs); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(buf.String(), "INFO: Logging from starlark with prefix") {
					t.Error("logger has unexpected log prefix and msg: ", buf.String())
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.test(t, test.kwargs)
		})
	}
}

func TestLogFuncScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		test   func(*testing.T, string)
	}{
		{
			name:   "logging with no  prefix",
			script: `log(msg="Logging from starlark")`,
			test: func(t *testing.T, script string) {
				exe := New()
				var buf bytes.Buffer
				logger := log.New(&buf, "", log.Lshortfile)
				exe.thread.SetLocal(identifiers.log, logger)
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(buf.String(), "Logging from starlark") {
					t.Error("logger has unexpected log msg: ", buf.String())
				}
			},
		},
		{
			name:   "logging with prefix",
			script: `log(prefix="INFO", msg="Logging from starlark with prefix")`,
			test: func(t *testing.T, script string) {
				exe := New()
				var buf bytes.Buffer
				logger := log.New(&buf, "", log.Lshortfile)
				exe.thread.SetLocal(identifiers.log, logger)
				if err := exe.Exec("test.star", strings.NewReader(script)); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(buf.String(), "INFO: Logging from starlark with prefix") {
					t.Error("logger has unexpected log prefix and msg: ", buf.String())
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.test(t, test.script)
		})
	}

}

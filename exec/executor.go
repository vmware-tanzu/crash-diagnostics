// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/crash-diagnostics/starlark"
)

type ArgMap map[string]string

func Execute(name string, source io.Reader, args ArgMap) error {
	star := starlark.New()

	if args != nil {
		starStruct, err := starlark.NewGoValue(args).ToStarlarkStruct("args")
		if err != nil {
			return err
		}

		star.AddPredeclared("args", starStruct)
	}

	err := star.Exec(name, source)
	if err != nil {
		err = errors.Wrap(err, "exec failed")
	}

	return err
}

func ExecuteFile(file *os.File, args ArgMap) error {
	return Execute(file.Name(), file, args)
}

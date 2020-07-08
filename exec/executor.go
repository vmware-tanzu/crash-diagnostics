// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"io"
	"os"

	"github.com/vmware-tanzu/crash-diagnostics/starlark"
)

type ArgMap map[string]string

func Execute(name string, source io.Reader, args ArgMap) error {
	star := starlark.New()

	if args != nil {
		starStruct, err := starlark.NewGoValue(args).ToStarlarkStruct()
		if err != nil {
			return err
		}

		star.AddPredeclared("args", starStruct)
	}

	if err := star.Exec(name, source); err != nil {
		return fmt.Errorf("exec failed: %s", err)
	}

	return nil
}

func ExecuteFile(file *os.File, args ArgMap) error {
	return Execute(file.Name(), file, args)
}

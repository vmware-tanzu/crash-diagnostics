package archive

import (
	"fmt"
	"os"

	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

func (c *cmd) Run(t *starlark.Thread, params Args) (Result, error) {
	if params.OutputFile == "" {
		params.OutputFile = DefaultBundleName
	}

	if len(params.SourcePaths) == 0 {
		return Result{Error: "no source path provided"}, nil
	}

	if err := archiver.Tar(params.OutputFile, params.SourcePaths...); err != nil {
		return Result{}, fmt.Errorf("%s failed: %w", FuncName, err)
	}

	info, err := os.Stat(params.OutputFile)
	if err != nil {
		return Result{}, fmt.Errorf("%s: stat failed: %w", FuncName, err)
	}

	return Result{Size: info.Size(), OutputFile: params.OutputFile}, nil
}

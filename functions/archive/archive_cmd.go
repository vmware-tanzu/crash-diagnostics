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

func (c *cmd) Run(t *starlark.Thread, params Args) Result {
	if params.OutputFile == "" {
		params.OutputFile = DefaultBundleName
	}

	if len(params.SourcePaths) == 0 {
		return Result{Error: "no source path provided"}
	}

	if err := archiver.Tar(params.OutputFile, params.SourcePaths...); err != nil {
		return Result{Error: fmt.Sprintf("%s failed: %s", Name, err)}
	}

	info, err := os.Stat(params.OutputFile)
	if err != nil {
		return Result{Error: fmt.Sprintf("%s: stat failed: %s", Name, err)}
	}

	return Result{Size: info.Size(), OutputFile: params.OutputFile}
}

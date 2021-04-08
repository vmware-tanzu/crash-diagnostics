package archive

import (
	"fmt"
	"os"

	"go.starlark.net/starlark"

	"github.com/vmware-tanzu/crash-diagnostics/archiver"
	"github.com/vmware-tanzu/crash-diagnostics/functions"
)

type cmd struct{}

func newCmd() *cmd {
	return new(cmd)
}

func (c *cmd) Run(t *starlark.Thread, p interface{}) (functions.CommandResult, error) {
	params, ok := p.(Params)
	if !ok {
		return nil, fmt.Errorf("unexpected param type: %T", p)
	}

	if len(params.OutputFile) == 0 {
		params.OutputFile = DefaultBundleName
	}

	if len(params.SourcePaths) == 0 {
		return functions.NewResult(params).AddError("no source path provided"), nil
	}

	if err := archiver.Tar(params.OutputFile, params.SourcePaths...); err != nil {
		return nil, fmt.Errorf("%s failed: %w", FuncName, err)
	}

	info, err := os.Stat(params.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("%s: stat failed: %w", FuncName, err)
	}

	params.Size = uint64(info.Size())

	return functions.NewResult(params), nil
}

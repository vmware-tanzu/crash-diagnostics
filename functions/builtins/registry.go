package builtins

import (
	"sync"

	"github.com/vmware-tanzu/crash-diagnostics/functions"
	"go.starlark.net/starlark"
)

var (
	mutex    sync.Mutex
	registry starlark.StringDict
)

func init() {
	registry = make(starlark.StringDict)
}

// Register registers a Starlark built-in function
func Register(name functions.FunctionName, builtin *starlark.Builtin) {
	mutex.Lock()
	defer mutex.Unlock()
	registry[string(name)] = builtin
}

func Registry() starlark.StringDict {
	return registry
}

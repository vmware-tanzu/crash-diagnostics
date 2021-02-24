package builtins

import (
	"sync"

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
func Register(name string, builtin *starlark.Builtin) {
	mutex.Lock()
	defer mutex.Unlock()
	registry[name] = builtin
}

func Registry() starlark.StringDict {
	return registry
}

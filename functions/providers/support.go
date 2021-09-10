package providers

import (
	"go.starlark.net/starlark"
)

func ResourcesFromThread(t *starlark.Thread) (Resources, bool) {
	if localVal := t.Local(ResourcesIdentifier); localVal != nil {
		resources, ok := localVal.(Resources)
		if !ok {
			return Resources{}, false
		}
		return resources, true
	}
	return Resources{}, false
}

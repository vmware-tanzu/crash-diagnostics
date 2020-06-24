package k8s

import (
	"go.starlark.net/starlark"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type SearchResult struct {
	ListKind             string
	ResourceName         string
	ResourceKind         string
	GroupVersionResource schema.GroupVersionResource
	List                 *unstructured.UnstructuredList
	Namespaced           bool
	Namespace            string
}

func (sr SearchResult) ToStarlarkValue() starlark.Value {
	var val starlark.Value
	return val
}

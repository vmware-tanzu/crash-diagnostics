package k8s

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type SearchParams struct {
	groups     []string
	kinds      []string
	namespaces []string
	versions   []string
	names      []string
	labels     []string
	containers []string
}

func (sp SearchParams) SetGroups(input []string) {
	sp.groups = input
}

func (sp SearchParams) SetKinds(input []string) {
	sp.kinds = input
}

func (sp SearchParams) SetNames(input []string) {
	sp.names = input
}

func (sp SearchParams) SetNamespaces(input []string) {
	sp.namespaces = input
}

func (sp SearchParams) SetVersions(input []string) {
	sp.versions = input
}

func (sp SearchParams) SetLabels(input []string) {
	sp.labels = input
}

func (sp SearchParams) SetContainers(input []string) {
	sp.containers = input
}

func (sp SearchParams) Groups() string {
	return strings.Join(sp.groups, " ")
}

func (sp SearchParams) Kinds() string {
	return strings.Join(sp.kinds, " ")
}

func (sp SearchParams) Names() string {
	return strings.Join(sp.names, " ")
}

func (sp SearchParams) Namespaces() string {
	return strings.Join(sp.namespaces, " ")
}

func (sp SearchParams) Versions() string {
	return strings.Join(sp.versions, " ")
}

func (sp SearchParams) Labels() string {
	return strings.Join(sp.labels, " ")
}

func (sp SearchParams) Containers() string {
	return strings.Join(sp.containers, " ")
}

// TODO: Change this to accept a string dictionary instead
func NewSearchParams(p *starlarkstruct.Struct) SearchParams {
	var (
		kinds      []string
		groups     []string
		names      []string
		namespaces []string
		versions   []string
		labels     []string
		containers []string
	)

	groups = parseStructAttr(p, "groups")
	kinds = parseStructAttr(p, "kinds")
	names = parseStructAttr(p, "names")
	namespaces = parseStructAttr(p, "namespaces")
	if len(namespaces) == 0 {
		namespaces = append(namespaces, "default")
	}
	versions = parseStructAttr(p, "versions")
	labels = parseStructAttr(p, "labels")
	containers = parseStructAttr(p, "containers")

	return SearchParams{
		kinds:      kinds,
		groups:     groups,
		names:      names,
		namespaces: namespaces,
		versions:   versions,
		labels:     labels,
		containers: containers,
	}
}

func parseStructAttr(p *starlarkstruct.Struct, attrName string) []string {
	values := make([]string, 0)

	attrVal, err := p.Attr(attrName)
	if err == nil {
		values, err = parse(attrVal)
		if err != nil {
			logrus.Errorf("error while parsing attr %s: %v", attrName, err)
		}
	}
	return values
}

func parse(inputValue starlark.Value) ([]string, error) {
	var values []string
	var err error

	switch inputValue.Type() {
	case "string":
		val, ok := inputValue.(starlark.String)
		if !ok {
			err = errors.Errorf("cannot process starlark value %s", inputValue.String())
			break
		}
		values = append(values, val.GoString())
	case "list":
		val, ok := inputValue.(*starlark.List)
		if !ok {
			err = errors.Errorf("cannot process starlark value %s", inputValue.String())
			break
		}
		iter := val.Iterate()
		defer iter.Done()
		var x starlark.Value
		for iter.Next(&x) {
			str, _ := x.(starlark.String)
			values = append(values, str.GoString())
		}
	default:
		err = errors.New("unknown input type for parse()")
	}

	return values, err
}

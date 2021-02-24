package functions

import (
	"go.starlark.net/starlark"
)

// ToStringSlice returns the elements in list as a []string
func ToStringSlice(list *starlark.List) []string {
	if list == nil {
		return []string{}
	}
	elems := make([]string, list.Len())
	for i := 0; i < list.Len(); i++ {
		if val, ok := list.Index(i).(starlark.String); ok {
			elems[i] = string(val)
		}
	}
	return elems
}

// ToIntSlice returns list elements as []int64
func ToIntSlice(list *starlark.List) []int64 {
	if list == nil {
		return []int64{}
	}
	elems := make([]int64, list.Len())
	for i := 0; i < list.Len(); i++ {
		if val, ok := list.Index(i).(starlark.Int).Int64(); ok {
			elems[i] = val
		}
	}
	return elems
}

// ToUintSlice returns list elements as []uint64
func ToUinttSlice(list *starlark.List) []uint64 {
	if list == nil {
		return []uint64{}
	}
	elems := make([]uint64, list.Len())
	for i := 0; i < list.Len(); i++ {
		if val, ok := list.Index(i).(starlark.Int).Uint64(); ok {
			elems[i] = val
		}
	}
	return elems
}

// ToFloatSlice returns list elements as []float64
func ToFloatSlice(list *starlark.List) []float64 {
	if list == nil {
		return []float64{}
	}
	elems := make([]float64, list.Len())
	for i := 0; i < list.Len(); i++ {
		if val, ok := list.Index(i).(starlark.Float); ok {
			elems[i] = float64(val)
		}
	}
	return elems
}

// ToBoolSlice returns list elements as []bool
func ToBoolSlice(list *starlark.List) []bool {
	if list == nil {
		return []bool{}
	}
	elems := make([]bool, list.Len())
	for i := 0; i < list.Len(); i++ {
		if val, ok := list.Index(i).(starlark.Bool); ok {
			elems[i] = bool(val)
		}
	}
	return elems
}
// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestGoValue_ToStringDict(t *testing.T) {
	tests := []struct {
		name  string
		goVal *GoValue
		eval  func(t *testing.T, goval *GoValue)
	}{
		{
			name:  "map[string]string",
			goVal: NewGoValue(map[string]string{"key0": "val0", "key1": "val1"}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[string]string)
				starVal, err := goval.ToStringDict()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected string
					if val, ok := starVal[k].(starlark.String); ok {
						expected = string(val)
					}
					if v != expected {
						t.Errorf("unexpected value not in starlark value: %s", k)
					}
				}
			},
		},
		{
			name:  "map[string]int",
			goVal: NewGoValue(map[string]int{"key0": 12, "key1": 14}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[string]int)
				starVal, err := goval.ToStringDict()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected int64
					if val, ok := starVal[k].(starlark.Int); ok {
						expected = val.BigInt().Int64()
					}
					if int64(v) != expected {
						t.Errorf("unexpected value not in starlark value: %s", k)
					}
				}
			},
		},
		{
			name:  "map[string]bool",
			goVal: NewGoValue(map[string]bool{"key0": false, "key1": true}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[string]bool)
				starVal, err := goval.ToStringDict()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected bool
					if val, ok := starVal[k].(starlark.Bool); ok {
						expected = bool(val)
					}
					if v != expected {
						t.Errorf("unexpected value not in starlark value: %s", k)
					}
				}
			},
		},
		{
			name:  "map[string][]string",
			goVal: NewGoValue(map[string][]string{"key0": []string{"hello", "goodbye"}, "key1": []string{"hi", "bye"}}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[string][]string)
				starVal, err := goval.ToStringDict()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected starlark.Tuple
					if val, ok := starVal[k].(starlark.Tuple); ok {
						expected = val
					}
					for i := range v {
						if v[i] != string(expected.Index(i).(starlark.String)) {
							t.Errorf("unexpected value not in starlark value: %s", expected.Index(i))
						}
					}
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.goVal)
		})
	}
}

func TestGoValue_ToDict(t *testing.T) {
	tests := []struct {
		name  string
		goVal *GoValue
		eval  func(t *testing.T, goval *GoValue)
	}{
		{
			name:  "map[string]string",
			goVal: NewGoValue(map[string]string{"key0": "val0", "key1": "val1"}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[string]string)
				dict, err := goval.ToDict()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected string

					if val, ok, err := dict.Get(starlark.String(k)); ok && err == nil {
						expected = string(val.(starlark.String))
					}
					if v != expected {
						t.Errorf("unexpected value not in starlark value: %s", v)
					}
				}
			},
		},
		{
			name:  "map[int]int",
			goVal: NewGoValue(map[int]int{0: 12, 10: 14}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[int]int)
				dict, err := goval.ToDict()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected int64
					if val, ok, err := dict.Get(starlark.MakeInt(k)); ok && err == nil {
						expected = val.(starlark.Int).BigInt().Int64()
					}
					if int64(v) != expected {
						t.Errorf("unexpected value not in starlark value: %v", v)
					}
				}
			},
		},
		{
			name:  "map[bool][]string",
			goVal: NewGoValue(map[bool][]string{true: []string{"hello", "goodbye"}, false: []string{"hi", "bye"}}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[bool][]string)
				dict, err := goval.ToDict()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected starlark.Tuple
					if val, ok, err := dict.Get(starlark.Bool(k)); ok && err == nil {
						expected = val.(starlark.Tuple)
					}
					for i := range v {
						if v[i] != string(expected.Index(i).(starlark.String)) {
							t.Errorf("unexpected value not in starlark value: %s", expected.Index(i))
						}
					}
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.goVal)
		})
	}
}

func TestGoValue_ToStruct(t *testing.T) {
	tests := []struct {
		name  string
		goVal *GoValue
		eval  func(t *testing.T, goval *GoValue)
	}{
		{
			name:  "map[string]string",
			goVal: NewGoValue(map[string]string{"key0": "val0", "key1": "val1"}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(map[string]string)
				starStruct, err := goval.ToStarlarkStruct()
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range actual {
					var expected string

					if val, err := starStruct.Attr(k); err == nil {
						expected = string(val.(starlark.String))
					}
					if v != expected {
						t.Errorf("unexpected value not in starlark value: %s", v)
					}
				}
			},
		},
		{
			name: "struct{string;int;bool}",
			goVal: NewGoValue(struct {
				Name  string
				Num   int
				Avail bool
			}{Name: "foo", Num: 10, Avail: true}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().(struct {
					Name  string
					Num   int
					Avail bool
				})
				starStruct, err := goval.ToStarlarkStruct()
				if err != nil {
					t.Fatal(err)
				}

				var attrName string
				if val, err := starStruct.Attr("Name"); err == nil {
					attrName = string(val.(starlark.String))
				}
				if actual.Name != attrName {
					t.Errorf("unexpected field value for 'name' starlark Struct : %s", attrName)
				}

				var attrNum int64
				if val, err := starStruct.Attr("Num"); err == nil {
					attrNum = val.(starlark.Int).BigInt().Int64()
				}
				if int64(actual.Num) != attrNum {
					t.Errorf("unexpected field value for 'num' starlark Struct : %d", attrNum)
				}

				var attrAvail bool
				if val, err := starStruct.Attr("Avail"); err == nil {
					attrAvail = bool(val.(starlark.Bool))
				}
				if actual.Avail != attrAvail {
					t.Errorf("unexpected field value for 'avail' starlark Struct : %t", attrAvail)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.goVal)
		})
	}
}

func TestGoValue_ToTuple(t *testing.T) {
	tests := []struct {
		name  string
		goVal *GoValue
		eval  func(t *testing.T, goval *GoValue)
	}{
		{
			name:  "[]string",
			goVal: NewGoValue([]string{"Hello", "World", "!"}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().([]string)
				tuple, err := goval.ToTuple()
				if err != nil {
					t.Fatal(err)
				}
				for i := range actual {
					var expected string
					if val, ok := tuple.Index(i).(starlark.String); ok {
						expected = string(val)
					}
					if actual[i] != expected {
						t.Errorf("unexpected value in starlark value: %s", actual[i])
					}
				}
			},
		},
		{
			name:  "[]bool",
			goVal: NewGoValue([]bool{true, true, false}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().([]bool)
				tuple, err := goval.ToTuple()
				if err != nil {
					t.Fatal(err)
				}
				for i := range actual {
					var expected bool
					if val, ok := tuple.Index(i).(starlark.Bool); ok {
						expected = bool(val)
					}
					if actual[i] != expected {
						t.Errorf("unexpected value in starlark value: %t", actual[i])
					}
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.goVal)
		})
	}
}

func TestGoValue_ToList(t *testing.T) {
	tests := []struct {
		name  string
		goVal *GoValue
		eval  func(t *testing.T, goval *GoValue)
	}{
		{
			name:  "[]string",
			goVal: NewGoValue([]string{"Hello", "World", "!"}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().([]string)
				list, err := goval.ToList()
				if err != nil {
					t.Fatal(err)
				}
				for i := range actual {
					var expected string
					if val, ok := list.Index(i).(starlark.String); ok {
						expected = string(val)
					}
					if actual[i] != expected {
						t.Errorf("unexpected value in starlark value: %s", actual[i])
					}
				}
			},
		},
		{
			name:  "[]bool",
			goVal: NewGoValue([]bool{true, true, false}),
			eval: func(t *testing.T, goval *GoValue) {
				actual := goval.Value().([]bool)
				list, err := goval.ToList()
				if err != nil {
					t.Fatal(err)
				}
				for i := range actual {
					var expected bool
					if val, ok := list.Index(i).(starlark.Bool); ok {
						expected = bool(val)
					}
					if actual[i] != expected {
						t.Errorf("unexpected value in starlark value: %t", actual[i])
					}
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.goVal)
		})
	}
}

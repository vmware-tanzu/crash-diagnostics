//// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
//// SPDX-License-Identifier: Apache-2.0

package typekit

import (
	"math"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestGoToStarlark(t *testing.T) {
	tests := []struct {
		name  string
		goVal interface{}
		eval  func(*testing.T, interface{})
	}{
		{
			name:  "Bool",
			goVal: true,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Bool
				if err := goToStarlark(true, &starval); err != nil {
					t.Fatal(err)
				}
				if starval != true {
					t.Errorf("unexpected bool value: %t", starval)
				}
			},
		},
		{
			name:  "Bool-Value",
			goVal: true,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(true, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.Truth() != true {
					t.Errorf("unexpected bool value: %t", starval)
				}
			},
		},
		{
			name:  "Int(Int64)",
			goVal: math.MaxInt32,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Int
				if err := goToStarlark(math.MaxInt32, &starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.Int64()
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val != math.MaxInt32 {
					t.Errorf("unexpected Int64 value: %d", val)
				}
			},
		},
		{
			name:  "Int(Uint64)",
			goVal: uint64(math.MaxUint64),
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Int
				if err := goToStarlark(uint64(math.MaxUint64), &starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.Uint64()
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val != math.MaxUint64 {
					t.Errorf("unexpected Uint64 value: %d", val)
				}
			},
		},
		{
			name:  "String",
			goVal: "Hello World!",
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.String
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if string(starval) != `Hello World!` {
					t.Errorf("unexpected string value: %s", starval)
				}
			},
		},
		{
			name:  "String-Value",
			goVal: "Hello World!",
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.String() != `"Hello World!"` {
					t.Errorf("unexpected string value: %s", starval)
				}
			},
		},
		{
			name:  "Tuple[string]",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				starval := make(starlark.Tuple, 2)
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}
				if starval.Index(1).String() != `"World!"` {
					t.Errorf("unexpected value: %s", starval.Index(1).String())
				}
			},
		},
		{
			name:  "Tuple[numeric]",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				starval := make(starlark.Tuple, 3)
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 3 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}

				intVal, _ := starval.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "Tuple[mix]",
			goVal: []interface{}{1, 2, 3, "Go!"},
			eval: func(t *testing.T, goVal interface{}) {
				starval := make(starlark.Tuple, 4)
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 4 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}
				strVal := starval.Index(3).String()
				if strVal != `"Go!"` {
					t.Errorf("Unexpected string element: %s", strVal)
				}
			},
		},
		{
			name:  "Tuple-Value",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				tuple := starval.(starlark.Tuple)
				if tuple.Len() != 3 {
					t.Errorf("unexpected tuple length %d", tuple.Len())
				}

				intVal, _ := tuple.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "Dict[string]string",
			goVal: map[string]string{"msg": "hello", "target": "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Dict
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected dict length %d", starval.Len())
				}
				val, _, err := starval.Get(starlark.String("target"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}

				if val.String() != `"world"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "Dict[string]int",
			goVal: map[string]int{"one": 12, "two": math.MaxInt8, "three": math.MaxInt64},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Dict
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 3 {
					t.Errorf("unexpected dict length %d", starval.Len())
				}
				val, _, err := starval.Get(starlark.String("three"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}
				if intVal, _ := val.(starlark.Int).Int64(); intVal != math.MaxInt64 {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "Dict-Value",
			goVal: map[string]int{"one": 12, "two": math.MaxInt8, "three": math.MaxInt64},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				dictVal := starval.(*starlark.Dict)
				if dictVal.Len() != 3 {
					t.Errorf("unexpected dict length %d", dictVal.Len())
				}
				val, _, err := dictVal.Get(starlark.String("three"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}
				if intVal, _ := val.(starlark.Int).Int64(); intVal != math.MaxInt64 {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "Struct-starlarkstruct",
			goVal: struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlarkstruct.Struct
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, err := starval.Attr("Msg")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name: "Struct-starlarkstruct-annotated",
			goVal: struct {
				Msg    string `name:"msg_field"`
				Target string
			}{
				Msg: "hello", Target: "world",
			},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlarkstruct.Struct
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, err := starval.Attr("msg_field")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "Struct-Value",
			goVal: struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				structVal := starval.(*starlarkstruct.Struct)
				val, err := structVal.Attr("Msg")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name: "Struct-Value-Annotated",
			goVal: struct {
				Msg    string
				Target string `name:"tgt"`
			}{
				Msg: "hello", Target: "world",
			},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				structVal := starval.(*starlarkstruct.Struct)
				val, err := structVal.Attr("tgt")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"world"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}
			},
		},
		{
			name:  "List[string]",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected list length %d", starval.Len())
				}
				if starval.Index(1).String() != `"World!"` {
					t.Errorf("unexpected list value: %s", starval.Index(1).String())
				}
			},
		},
		{
			name:  "List[numeric]",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 3 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}

				intVal, _ := starval.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "Tuple[mix]",
			goVal: []interface{}{1, 2, 3, "Go!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 4 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}
				strVal := starval.Index(3).String()
				if strVal != `"Go!"` {
					t.Errorf("Unexpected string element: %s", strVal)
				}
			},
		},
		{
			name:  "List-Value",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				list := starval.(*starlark.List)
				if list.Len() != 3 {
					t.Errorf("unexpected tuple length %d", list.Len())
				}

				intVal, _ := list.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "Set[string]",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Set
				if err := Go(goVal).StarlarkSet(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected set length %d", starval.Len())
				}
				iter := starval.Iterate()
				var val starlark.Value
				for iter.Next(&val) {
					if val.String() == `"World!"` {
						iter.Done()
						return
					}
				}
				t.Errorf("Stararlk set value not found")
			},
		},
		{
			name:  "Set[mix]",
			goVal: []interface{}{1, 2, 3, "Go!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 4 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}
				iter := starval.Iterate()
				var val starlark.Value
				for iter.Next(&val) {
					if val.String() == `"Go!"` {
						iter.Done()
						return
					}
				}
				t.Errorf("Stararlk set value not found")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.goVal)
		})
	}
}

package typekit

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestStarlarkToGo(t *testing.T) {
	tests := []struct {
		name    string
		starVal starlark.Value
		eval    func(*testing.T, starlark.Value)
	}{
		{
			name:    "bool",
			starVal: starlark.Bool(true),
			eval: func(t *testing.T, val starlark.Value) {
				var boolVar bool
				err := starlarkToGo(val, reflect.ValueOf(&boolVar).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if !boolVar {
					t.Fatalf("unexpected bool value: %t", boolVar)
				}
			},
		},
		{
			name:    "int32",
			starVal: starlark.MakeInt(math.MaxInt32),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar int64
				err := starlarkToGo(val, reflect.ValueOf(&intVar).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxInt32 {
					t.Fatalf("unexpected int32 value: %d", intVar)
				}
			},
		},
		{
			name:    "int64",
			starVal: starlark.MakeInt64(math.MaxInt64),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar int64
				err := starlarkToGo(val, reflect.ValueOf(&intVar).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxInt64 {
					t.Fatalf("unexpected int32 value: %d", intVar)
				}
			},
		},
		{
			name:    "uint64",
			starVal: starlark.MakeUint64(math.MaxUint64),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar uint64
				err := starlarkToGo(val, reflect.ValueOf(&intVar).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxUint64 {
					t.Fatalf("unexpected int32 value: %d", intVar)
				}
			},
		},
		{
			name:    "float32",
			starVal: starlark.Float(math.MaxFloat32),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar float64
				err := starlarkToGo(val, reflect.ValueOf(&floatVar).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if floatVar != math.MaxFloat32 {
					t.Fatalf("unexpected float64 value: %f", floatVar)
				}
			},
		},
		{
			name:    "float64",
			starVal: starlark.Float(math.MaxFloat64),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar float64
				err := starlarkToGo(val, reflect.ValueOf(&floatVar).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if floatVar != math.MaxFloat64 {
					t.Fatalf("unexpected float64 value: %f", floatVar)
				}
			},
		},
		{
			name:    "string",
			starVal: starlark.String("Hello World!"),
			eval: func(t *testing.T, val starlark.Value) {
				var strVar string
				err := starlarkToGo(val, reflect.ValueOf(&strVar).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if strVar != "Hello World!" {
					t.Fatalf("unexpected string value: %s", strVar)
				}
			},
		},
		{
			name:    "list-string",
			starVal: starlark.NewList([]starlark.Value{starlark.String("Hello"), starlark.String("World!")}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]string, 0)
				err := starlarkToGo(val, reflect.ValueOf(&slice).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if "Hello World!" != strings.Join(slice, " ") {
					t.Fatalf("unexpected string value: %v", slice)
				}
			},
		},
		{
			name:    "list-numbers",
			starVal: starlark.NewList([]starlark.Value{starlark.MakeInt64(math.MaxInt64), starlark.MakeInt(math.MaxInt32)}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]int64, 0)
				err := starlarkToGo(val, reflect.ValueOf(&slice).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0] != math.MaxInt64 {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1] != math.MaxInt32 {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
			},
		},
		{
			name:    "list-mixed",
			starVal: starlark.NewList([]starlark.Value{starlark.String("HelloWorld!"), starlark.MakeInt(math.MaxInt32)}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := starlarkToGo(val, reflect.ValueOf(&slice).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name:    "tuple-mixed",
			starVal: starlark.Tuple([]starlark.Value{starlark.String("HelloWorld!"), starlark.MakeInt(math.MaxInt32)}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := starlarkToGo(val, reflect.ValueOf(&slice).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name: "dict[string]string",
			starVal: func() *starlark.Dict {
				dict := starlark.NewDict(2)
				dict.SetKey(starlark.String("msg0"), starlark.String("Hello"))
				dict.SetKey(starlark.String("msg1"), starlark.String("World!"))
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]string)
				err := starlarkToGo(val, reflect.ValueOf(&gomap).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gomap["msg0"] != "Hello" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
				if gomap["msg1"] != "World!" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
			},
		},
		{
			name: "dict[string]int",
			starVal: func() *starlark.Dict {
				dict := starlark.NewDict(2)
				dict.SetKey(starlark.String("msg0"), starlark.MakeInt(math.MaxInt32))
				dict.SetKey(starlark.String("msg1"), starlark.MakeInt64(math.MaxInt64))
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]int64)
				err := starlarkToGo(val, reflect.ValueOf(&gomap).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gomap["msg0"] != math.MaxInt32 {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
				if gomap["msg1"] != math.MaxInt64 {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
			},
		},
		{
			name: "set-string",
			starVal: func() *starlark.Set {
				set := starlark.NewSet(2)
				set.Insert(starlark.String("HelloWorld!"))
				set.Insert(starlark.MakeInt(math.MaxInt32))
				return set
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := starlarkToGo(val, reflect.ValueOf(&slice).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name: "set-mixed",
			starVal: func() *starlark.Set {
				set := starlark.NewSet(2)
				set.Insert(starlark.String("HelloWorld!"))
				set.Insert(starlark.MakeInt(math.MaxInt32))
				return set
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := starlarkToGo(val, reflect.ValueOf(&slice).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name: "struct",
			starVal: func() *starlarkstruct.Struct {
				dict := starlark.StringDict{
					"msg0": starlark.String("Hello"),
					"msg1": starlark.String("World!"),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				var gostruct struct{ Msg0, Msg1 string }
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Msg0 != "Hello" {
					t.Fatalf("unexpected map[msg] value: %v", gostruct.Msg0)
				}
				if gostruct.Msg1 != "World!" {
					t.Fatalf("unexpected map[msg] value: %v", gostruct.Msg1)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.starVal)
		})
	}
}

func TestStarGo(t *testing.T) {
	tests := []struct {
		name    string
		starVal starlark.Value
		eval    func(*testing.T, starlark.Value)
	}{
		{
			name:    "bool",
			starVal: starlark.Bool(true),
			eval: func(t *testing.T, val starlark.Value) {
				var boolVar bool
				if err := Starlark(val).Go(&boolVar); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if !boolVar {
					t.Fatalf("unexpected bool value: %t", boolVar)
				}
			},
		},

		{
			name:    "int64",
			starVal: starlark.MakeInt64(math.MaxInt64),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar int64
				if err := Starlark(val).Go(&intVar); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxInt64 {
					t.Fatalf("unexpected int32 value: %d", intVar)
				}
			},
		},
		{
			name:    "float64",
			starVal: starlark.Float(math.MaxFloat64),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar float64
				if err := Starlark(val).Go(&floatVar); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if floatVar != math.MaxFloat64 {
					t.Fatalf("unexpected float64 value: %f", floatVar)
				}
			},
		},
		{
			name:    "list-string",
			starVal: starlark.NewList([]starlark.Value{starlark.String("Hello"), starlark.String("World!")}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]string, 0)
				if err := Starlark(val).Go(&slice); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if "Hello World!" != strings.Join(slice, " ") {
					t.Fatalf("unexpected string value: %v", slice)
				}
			},
		},
		{
			name: "dict[string]string",
			starVal: func() *starlark.Dict {
				dict := starlark.NewDict(2)
				dict.SetKey(starlark.String("msg0"), starlark.String("Hello"))
				dict.SetKey(starlark.String("msg1"), starlark.String("World!"))
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]string)
				if err := Starlark(val).Go(&gomap); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gomap["msg0"] != "Hello" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
				if gomap["msg1"] != "World!" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.starVal)
		})
	}
}

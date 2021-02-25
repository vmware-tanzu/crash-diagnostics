package typekit

import (
	"reflect"
	"strings"
	"testing"

	"go.starlark.net/starlark"
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
		//	{
		//		name:    "int32",
		//		starVal: starlark.MakeInt(math.MaxInt32),
		//		eval: func(t *testing.T, val starlark.Value) {
		//			var intVar int64
		//			err := starlarkToGo(val, &intVar)
		//			if err != nil {
		//				t.Fatalf("failed to convert starlark to go value: %s", err)
		//			}
		//			if intVar != math.MaxInt32 {
		//				t.Fatalf("unexpected int32 value: %d", intVar)
		//			}
		//		},
		//	},
		//	{
		//		name:    "int64",
		//		starVal: starlark.MakeInt64(math.MaxInt64),
		//		eval: func(t *testing.T, val starlark.Value) {
		//			var intVar int64
		//			err := starlarkToGo(val, &intVar)
		//			if err != nil {
		//				t.Fatalf("failed to convert starlark to go value: %s", err)
		//			}
		//			if intVar != math.MaxInt64 {
		//				t.Fatalf("unexpected int32 value: %d", intVar)
		//			}
		//		},
		//	},
		//	{
		//		name:    "uint64",
		//		starVal: starlark.MakeUint64(math.MaxUint64),
		//		eval: func(t *testing.T, val starlark.Value) {
		//			var intVar uint64
		//			err := starlarkToGo(val, &intVar)
		//			if err != nil {
		//				t.Fatalf("failed to convert starlark to go value: %s", err)
		//			}
		//			if intVar != math.MaxUint64 {
		//				t.Fatalf("unexpected int32 value: %d", intVar)
		//			}
		//		},
		//	},
		//	{
		//		name:    "float32",
		//		starVal: starlark.Float(math.MaxFloat32),
		//		eval: func(t *testing.T, val starlark.Value) {
		//			var floatVar float64
		//			err := starlarkToGo(val, &floatVar)
		//			if err != nil {
		//				t.Fatalf("failed to convert starlark to go value: %s", err)
		//			}
		//			if floatVar != math.MaxFloat32 {
		//				t.Fatalf("unexpected float64 value: %f", floatVar)
		//			}
		//		},
		//	},
		//	{
		//		name:    "float64",
		//		starVal: starlark.Float(math.MaxFloat64),
		//		eval: func(t *testing.T, val starlark.Value) {
		//			var floatVar float64
		//			err := starlarkToGo(val, &floatVar)
		//			if err != nil {
		//				t.Fatalf("failed to convert starlark to go value: %s", err)
		//			}
		//			if floatVar != math.MaxFloat64 {
		//				t.Fatalf("unexpected float64 value: %f", floatVar)
		//			}
		//		},
		//	},
		//	{
		//		name:    "string",
		//		starVal: starlark.String("Hello World!"),
		//		eval: func(t *testing.T, val starlark.Value) {
		//			var strVar string
		//			err := starlarkToGo(val, &strVar)
		//			if err != nil {
		//				t.Fatalf("failed to convert starlark to go value: %s", err)
		//			}
		//			if strVar != "Hello World!" {
		//				t.Fatalf("unexpected string value: %s", strVar)
		//			}
		//		},
		//	},
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.starVal)
		})
	}
}

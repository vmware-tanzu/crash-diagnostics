// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package typekit

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
)

type StarValue struct {
	val starlark.Value
}

func Starval(val starlark.Value) *StarValue {
	return &StarValue{val: val}
}

func (v *StarValue) Value() starlark.Value {
	return v.val
}

// ToStringSlice returns the elements in list as a []string
func (v *StarValue) ToStringSlice() ([]string, error) {
	if err := assertVal(v.val); err != nil {
		return nil, fmt.Errorf("StarValue: %s", err)
	}

	list, ok := v.val.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("Starval.ToStringSlice: type %T conversion failed", v.val)
	}

	elems := make([]string, list.Len())
	for i := 0; i < list.Len(); i++ {
		if val, ok := list.Index(i).(starlark.String); ok {
			elems[i] = string(val)
		}
	}
	return elems, nil
}

// starlarkToGo translates starlark.Value val to the provided Go value goval
// using the following type mapping:
//
//      starlark.Bool   	-- bool
//      starlark.Int    	-- int64 or uint64
//      starlark.Float  	-- float64
//      starlark.String 	-- string
//      *starlark.List  	-- []T
//      starlark.Tuple  	-- []T
//      *starlark.Dict  	-- map[K]T
//      *starlark.Set   	-- []T

func starlarkToGo(val starlark.Value, goval reflect.Value) error {
	gotype := goval.Type()

	var starval reflect.Value
	switch val.Type() {
	case "bool":
		if gotype.Kind() != reflect.Bool {
			return fmt.Errorf("target type must be bool")
		}
		starval = reflect.ValueOf(bool(val.Truth()))

	case "int":
		if gotype.Kind() != reflect.Int64 && gotype.Kind() != reflect.Uint64 {
			return fmt.Errorf("target type must be int64 or uint64")
		}
		intVal, ok := val.(starlark.Int)
		if !ok {
			return fmt.Errorf("failed to assert %v as starlark.Int", val)
		}

		bigInt := intVal.BigInt()
		switch {
		case bigInt.IsInt64():
			starval = reflect.ValueOf(bigInt.Int64())
		case bigInt.IsUint64():
			starval = reflect.ValueOf(bigInt.Uint64())
		default:
			return fmt.Errorf("unsupported starlark.Int type")
		}

	case "float":
		if gotype.Kind() != reflect.Float64 {
			return fmt.Errorf("target type must be float64")
		}
		floatVal, ok := val.(starlark.Float)
		if !ok {
			return fmt.Errorf("failed to assert %v as starlark.Float", val)
		}
		starval = reflect.ValueOf(float64(floatVal))

	case "string":
		if gotype.Kind() != reflect.String {
			return fmt.Errorf("target type must be string")
		}
		strVal, ok := val.(starlark.String)
		if !ok {
			return fmt.Errorf("failed to assert %v as starlark.String", val)
		}
		starval = reflect.ValueOf(string(strVal))

	case "list":
		if gotype.Kind() != reflect.Slice && gotype.Kind() != reflect.Array {
			return fmt.Errorf("target type must be slice or array")
		}
		listVal, ok := val.(*starlark.List)
		if !ok {
			return fmt.Errorf("failed to assert %v as *starlark.List", val)
		}
		goval.Set(reflect.MakeSlice(gotype, listVal.Len(), listVal.Len()))
		for i := 0; i < listVal.Len(); i++ {
			if err := starlarkToGo(listVal.Index(i), goval.Index(i)); err != nil {
				return err
			}
		}
		return nil
	}
	if gotype.AssignableTo(starval.Type()) {
		goval.Set(starval.Convert(gotype))
		return nil
	}

	//	slice := make([]interface{}, listVal.Len())
	//	for i := 0; i < listVal.Len(); i++ {
	//		goVal, err := starlarkToGo(listVal.Index(i))
	//		if err != nil {
	//			return nil, fmt.Errorf("failed to convert list item to Go: %s", err)
	//		}
	//		slice[i] = goVal
	//	}
	//
	//	return slice, nil

	//case "dict":
	//case "set":
	//case "function":
	//default:
	//	return fmt.Errorf("unable to convert Starlark type %s to Go type", val.Type())
	//}
	return fmt.Errorf("failed conversion")
}

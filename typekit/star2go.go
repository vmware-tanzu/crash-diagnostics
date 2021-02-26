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

// Star wraps a Starlark value
func Star(val starlark.Value) *StarValue {
	return &StarValue{val: val}
}

func (v *StarValue) Value() starlark.Value {
	return v.val
}

// ToGo converts Starlark value in StarValue into the Go
// value specified by pointer to value goPtr.
// Example:
//
//    var msg string
//    Star(starlark.String("Hello")).Go(&msg)
//
// This method supports the following type map from Starlark to Go types:
//
//      starlark.Bool   	-- bool
//      starlark.Int    	-- int64 or uint64
//      starlark.Float  	-- float64
//      starlark.String 	-- string
//      *starlark.List  	-- []T
//      starlark.Tuple  	-- []T
//      *starlark.Dict  	-- map[K]T
//      *starlark.Set   	-- []T

func (v *StarValue) ToGo(goPtr interface{}) error {
	goval := reflect.ValueOf(goPtr)
	gotype := goval.Type()
	if gotype.Kind() != reflect.Ptr || goval.IsNil() {
		return fmt.Errorf("requires non-nil pointer, got %v", gotype)
	}

	return starlarkToGo(v.val, goval.Elem())
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

func starlarkToGo(srcVal starlark.Value, goval reflect.Value) error {
	gotype := goval.Type()

	var starval reflect.Value
	switch srcVal.Type() {
	case "bool":
		if gotype.Kind() != reflect.Bool && gotype.Kind() != reflect.Interface {
			return fmt.Errorf("target type must be bool")
		}
		starval = reflect.ValueOf(bool(srcVal.Truth()))

	case "int":
		if gotype.Kind() != reflect.Int64 && gotype.Kind() != reflect.Uint64 && gotype.Kind() != reflect.Interface {
			return fmt.Errorf("target type must be int64 or uint64")
		}
		intVal, ok := srcVal.(starlark.Int)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Int", srcVal)
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
		if gotype.Kind() != reflect.Float64 && gotype.Kind() != reflect.Interface {
			return fmt.Errorf("target type must be float64")
		}
		floatVal, ok := srcVal.(starlark.Float)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Float", srcVal)
		}
		starval = reflect.ValueOf(float64(floatVal))

	case "string":
		if gotype.Kind() != reflect.String && gotype.Kind() != reflect.Interface {
			return fmt.Errorf("target type must be string or interface{}")
		}
		strVal, ok := srcVal.(starlark.String)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.String", srcVal)
		}
		starval = reflect.ValueOf(string(strVal))

	case "list":
		if gotype.Kind() != reflect.Slice && gotype.Kind() != reflect.Array {
			return fmt.Errorf("target type must be slice or array")
		}
		listVal, ok := srcVal.(*starlark.List)
		if !ok {
			return fmt.Errorf("failed to assert %T as *starlark.List", srcVal)
		}
		goval.Set(reflect.MakeSlice(gotype, listVal.Len(), listVal.Len()))
		for i := 0; i < listVal.Len(); i++ {
			if err := starlarkToGo(listVal.Index(i), goval.Index(i)); err != nil {
				return err
			}
		}
		return nil

	case "tuple":
		if gotype.Kind() != reflect.Slice && gotype.Kind() != reflect.Array {
			return fmt.Errorf("target type must be slice or array")
		}
		tupVal, ok := srcVal.(starlark.Tuple)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Tuple", srcVal)
		}
		goval.Set(reflect.MakeSlice(gotype, tupVal.Len(), tupVal.Len()))
		for i := 0; i < tupVal.Len(); i++ {
			if err := starlarkToGo(tupVal.Index(i), goval.Index(i)); err != nil {
				return err
			}
		}
		return nil

	case "dict":
		if gotype.Kind() != reflect.Map {
			return fmt.Errorf("target type must be map")
		}
		dict, ok := srcVal.(*starlark.Dict)
		if !ok {
			return fmt.Errorf("failed to assert %T as *starlark.Dict", srcVal)
		}
		goval.Set(reflect.MakeMap(gotype))
		for _, dictKey := range dict.Keys() {
			dictVal, ok, err := dict.Get(dictKey)
			if err != nil {
				return fmt.Errorf("starlark.Dict.Get failed: %s", err)
			}
			if !ok {
				continue
			}

			// convert starlark key to Go value
			goMapKey := reflect.New(gotype.Key()).Elem()
			if err := starlarkToGo(dictKey, goMapKey); err != nil {
				return err
			}

			// convert starlark dict value to Go value
			goMapVal := reflect.New(gotype.Elem()).Elem()
			if err := starlarkToGo(dictVal, goMapVal); err != nil {
				return err
			}

			// store map value
			goval.SetMapIndex(goMapKey, goMapVal)
		}
		return nil

	case "set":
		if gotype.Kind() != reflect.Slice && gotype.Kind() != reflect.Array {
			return fmt.Errorf("target type must be slice or array")
		}
		setVal, ok := srcVal.(*starlark.Set)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Tuple", srcVal)
		}
		goval.Set(reflect.MakeSlice(gotype, setVal.Len(), setVal.Len()))
		var setItem starlark.Value
		iter := setVal.Iterate()
		i := 0
		for iter.Next(&setItem) {
			if err := starlarkToGo(setItem, goval.Index(i)); err != nil {
				return err
			}
			i++
		}
		return nil
	}

	startype := starval.Type()
	if startype.ConvertibleTo(gotype) {
		goval.Set(starval.Convert(gotype))
		return nil
	}

	return fmt.Errorf("failed conversion")
}

//// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
//// SPDX-License-Identifier: Apache-2.0

package typekit

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
)

// KwargsToGo converts a slice of Starlark kwargs (keyword args) value
// to a Go struct. The function uses annotated fields on the struct to describe
// the keyword argument mapping as:
//
//  type Param struct {
//      OutputFile  string   `name:"output_file" optional:"true"`
//      SourcePaths []string `name:"output_path"`
//  }
//
// The previous struct can be mapped to the following keyword args:
//
// []starlark.Tuple{
//     {starlark.String("output_file"), starlark.String("/tmp/out.tar.gz")},
//	   {starlark.String("source_paths"), starlark.NewList([]starlark.Value{starlark.String("/tmp/crashd")})},
// }
//
// Supported annotation: `name:"arg_name" optional:"true|false" (default false)`
func KwargsToGo(kwargs []starlark.Tuple, goStructPtr interface{}) error {
	goval := reflect.ValueOf(goStructPtr)
	gotype := goval.Type()
	if gotype.Kind() != reflect.Ptr || goval.IsNil() {
		return fmt.Errorf("kwargs expects a non-nil pointer to a struct, got %v", gotype.Kind())
	}
	return kwargsToGo(kwargs, goval.Elem())
}

func kwargsToGo(kwargs []starlark.Tuple, goval reflect.Value) error {
	gotype := goval.Type()
	if gotype.Kind() != reflect.Struct {
		return fmt.Errorf("kwargs requires non-nil pointer to struct, got: %v", gotype)
	}

	if !goval.IsValid() {
		goval.Set(reflect.Zero(goval.Type()))
	}

	for i := 0; i < goval.NumField(); i++ {
		field := gotype.Field(i)

		argName, ok := field.Tag.Lookup("name")
		if !ok {
			continue
		}

		// get arg from keyword args (use either tag or field name)
		kwarg, err := getKwarg(kwargs, argName, field.Name)
		if err != nil {
			return err
		}

		// is arg marked optional? By default args are optional=false
		// arg is optional if it is explicitly marked with "true" or "yes"
		argOptional, _ := field.Tag.Lookup("optional")
		switch argOptional {
		case "true", "yes":
		default:
			if kwarg == starlark.None {
				return fmt.Errorf("argument '%s' is required", argName)
			}
		}

		// set field value if not None
		if kwarg != starlark.None {
			fieldVal := reflect.New(field.Type).Elem()
			if err := starlarkToGo(kwarg, fieldVal); err != nil {
				return err
			}
			goval.FieldByName(field.Name).Set(fieldVal)
		}
	}

	return nil
}

func getKwarg(kwargs []starlark.Tuple, argName, defaultName string) (starlark.Value, error) {
	for _, kwarg := range kwargs {
		arg, ok := kwarg.Index(0).(starlark.String)
		if !ok {
			return nil, fmt.Errorf("keyword arg name is not a string")
		}
		if string(arg) == argName || string(arg) == defaultName {
			return kwarg.Index(1), nil
		}
	}
	return starlark.None, nil
}

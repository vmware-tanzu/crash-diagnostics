package typekit

import (
	"fmt"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestStarKwargs2Go(t *testing.T) {
	tests := []struct {
		name   string
		kwargs []starlark.Tuple
		eval   func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name:   "missing explicit optional arg",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a" optional:"true"`
				}{}
				if err := KwargsToGo(kwargs, &val); err != nil {
					t.Fatal(err)
				}
				if val.A != "" {
					t.Error("expecting empty value")
				}
			},
		},
		{
			name:   "missing explicit required arg",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a" optional:"false"`
				}{}
				if err := KwargsToGo(kwargs, &val); err == nil {
					t.Fatal("should fail due to missing arg")
				}
			},
		},
		{
			name:   "missing default required arg",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a"`
				}{}
				if err := KwargsToGo(kwargs, &val); err == nil {
					t.Fatal("should fail due to missing arg")
				}
			},
		},
		{
			name: "all required",
			kwargs: []starlark.Tuple{
				{starlark.String("a"), starlark.String("hello")},
				{starlark.String("b"), starlark.MakeInt(32)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a"`
					B int64  `name:"b"`
				}{}
				if err := KwargsToGo(kwargs, &val); err != nil {
					t.Fatal(err)
				}
				if val.A != "hello" {
					t.Errorf("unexpected value: %s", val.A)
				}
				if val.B != 32 {
					t.Errorf("unexpected value: %d", val.B)
				}
			},
		},
		{
			name: "all optional",
			kwargs: []starlark.Tuple{
				{starlark.String("a"), starlark.String("hello")},
				{starlark.String("b"), starlark.MakeInt(32)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a" optional:"true"`
					B int64  `name:"b" optional:"true"`
				}{}
				if err := KwargsToGo(kwargs, &val); err != nil {
					t.Fatal(err)
				}
				if val.A != "hello" {
					t.Errorf("unexpected value: %s", val.A)
				}
				if val.B != 32 {
					t.Errorf("unexpected value: %d", val.B)
				}
			},
		},
		{
			name: "optional and required",
			kwargs: []starlark.Tuple{
				{starlark.String("a"), starlark.String("hello")},
				{starlark.String("b"), starlark.MakeInt(32)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a" optional:"true"`
					B int64  `name:"b" optional:"false"`
				}{}
				if err := KwargsToGo(kwargs, &val); err != nil {
					t.Fatal(err)
				}
				if val.A != "hello" {
					t.Errorf("unexpected value: %s", val.A)
				}
				if val.B != 32 {
					t.Errorf("unexpected value: %d", val.B)
				}
			},
		},
		{
			name: "implicit require and explicit optional",
			kwargs: []starlark.Tuple{
				{starlark.String("a"), starlark.String("hello")},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a"`
					B int64  `name:"b" optional:"true"`
				}{}
				if err := KwargsToGo(kwargs, &val); err != nil {
					t.Fatal(err)
				}
				if val.A != "hello" {
					t.Errorf("unexpected value: %s", val.A)
				}
				if val.B != 0 {
					t.Errorf("unexpected value: %d", val.B)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.kwargs)
		})
	}
}

func TestGoToKwargToGo(t *testing.T) {
	type Desc struct {
		Name string
	}
	type inarg struct {
		Name  Desc  `name:"name"`
		Count int64 `name:"count"`
	}

	// kwarg -> Go
	var arg inarg
	err := KwargsToGo(
		[]starlark.Tuple{
			{
				starlark.String("name"),
				starlarkstruct.FromStringDict(
					starlarkstruct.Default,
					starlark.StringDict{"name": starlark.String("Hello")},
				),
			},
			{starlark.String("count"), starlark.MakeInt(266)},
		},
		&arg,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Go -> starlark struct
	starOut := new(starlarkstruct.Struct)
	if err := Go(arg).Starlark(starOut); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %v", err))
	}

	// starlark struct -> Go
	var arg2 inarg
	if err := Starlark(starOut).Go(&arg2); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %w", err))
	}

	if arg2.Name.Name != "Hello" {
		t.Errorf("Unexpected value: %s", arg2.Name.Name)
	}
}

func TestGoToKwargToGo2(t *testing.T) {
	type Desc struct {
		Name string
	}
	type inarg struct {
		Name  Desc  `name:"name"`
		Count int64 `name:"count"`
		B     string
	}

	// kwarg -> Go
	var arg inarg
	desc := Desc{Name: "World"}
	descArg := new(starlarkstruct.Struct)
	if err := Go(desc).Starlark(descArg); err != nil {
		t.Fatal(err)
	}
	err := KwargsToGo(
		[]starlark.Tuple{
			{starlark.String("name"), descArg},
			{starlark.String("count"), starlark.MakeInt(266)},
		},
		&arg,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Go -> starlark struct
	starOut := new(starlarkstruct.Struct)
	if err := Go(arg).Starlark(starOut); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %v", err))
	}

	// starlark struct -> Go
	var arg2 inarg
	if err := Starlark(starOut).Go(&arg2); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %w", err))
	}

	if arg2.Name.Name != "World" {
		t.Errorf("Unexpected value: %s", arg2.Name.Name)
	}
}

package typekit

import (
	"testing"

	"go.starlark.net/starlark"
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
					A string `arg:"a" optional:"true"`
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
					A string `arg:"a" optional:"false"`
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
					A string `arg:"a"`
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
					A string `arg:"a"`
					B int64  `arg:"b"`
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
					A string `arg:"a" optional:"true"`
					B int64  `arg:"b" optional:"true"`
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
					A string `arg:"a" optional:"true"`
					B int64  `arg:"b" optional:"false"`
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
					A string `arg:"a"`
					B int64  `arg:"b" optional:"true"`
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

package lang

import (
	"reflect"
	"testing"
)

func TestMatch(t *testing.T) {
	type T struct {
		name    string
		in      any
		matched bool
		out     any
	}
	cases := []struct {
		name string
		m    Matcher[any]
		xs   []T
	}{
		{
			name: "string",
			m: SyntaxRule(
				Case[any](
					"hello",
					func(w W, a ...any) any {
						return nil
					},
				),
			),
			xs: []T{
				{
					name:    "hello",
					in:      "hello",
					matched: true,
					out:     nil,
				},
				{
					name:    "world",
					in:      "world",
					matched: false,
					out:     nil,
				},
			},
		},
		{
			name: "int",
			m: SyntaxRule(
				Case[any](
					42,
					func(w W, a ...any) any {
						return nil
					},
				),
			),
			xs: []T{
				{
					name:    "42",
					in:      42,
					matched: true,
					out:     nil,
				},
				{
					name:    "100",
					in:      100,
					matched: false,
					out:     nil,
				},
			},
		},
		{
			name: "BuiltPattern",
			m: SyntaxRule(
				Case[any](
					BuiltPattern(func(n Node, w W) bool {
						w["x"] = 42
						return true
					}),
					func(w W, a ...any) any {
						return w["x"].(int)
					},
				),
			),
			xs: []T{
				{
					name:    "42",
					in:      42,
					matched: true,
					out:     42,
				},
			},
		},
		{
			name: "Slot Any",
			m: SyntaxRule(
				Case[any](
					S("a"),
					func(w W, a ...any) any {
						return w["a"]
					},
				),
			),
			xs: []T{
				{
					name:    "42",
					in:      42,
					matched: true,
					out:     42,
				},
			},
		},
		{
			name: "Slot With Nested BuiltPattern",
			m: SyntaxRule(
				Case[any](
					S("a", BuiltPattern(func(n Node, w W) bool {
						return n.(int) == 42
					})),
					func(w W, a ...any) any {
						return w["a"]
					},
				),
			),
			xs: []T{
				{
					name:    "42",
					in:      42,
					matched: true,
					out:     42,
				},
				{
					name:    "100",
					in:      100,
					matched: false,
					out:     nil,
				},
			},
		},

		{
			name: "Slot With Nested List",
			m: SyntaxRule(
				Case[any](
					S("a", L(1, 2, 3, 4)),
					func(w W, a ...any) any {
						return w["a"]
					},
				),
			),
			xs: []T{
				{
					name:    "[1, 2, 3, 4]",
					in:      []any{1, 2, 3, 4},
					matched: true,
					out:     []any{1, 2, 3, 4},
				},
				{
					name:    "[1, 2, 3]",
					in:      []any{1, 2, 3},
					matched: false,
					out:     nil,
				},
			},
		},
		// todo empty list
		{
			name: "Trival List",
			m: SyntaxRule(
				Case[any](
					L(1, "2", 3, "4"),
					func(w W, a ...any) any {
						return nil
					},
				),
			),
			xs: []T{
				{
					name:    "[1, '2', 3, '4']",
					in:      []any{1, "2", 3, "4"},
					matched: true,
					out:     nil,
				},
				{
					name:    "[1, 2, 3, 4]",
					in:      []any{1, 2, 3, 4},
					matched: false,
					out:     nil,
				},
			},
		},
		{
			name: "Nested List",
			m: SyntaxRule(
				Case[any](
					L(1, L(2, 3), 4),
					func(w W, a ...any) any {
						return nil
					},
				),
			),
			xs: []T{
				{
					name:    "[1, [2, 3], 4]",
					in:      []any{1, []any{2, 3}, 4},
					matched: true,
					out:     nil,
				},
				{
					name:    "[1, 2, 3, 4]",
					in:      []any{1, 2, 3, 4},
					matched: false,
					out:     nil,
				},
			},
		},
		{
			name: "Nested Slot List",
			m: SyntaxRule(
				Case[any](
					L(1, L(2, S("a")), 4),
					func(w W, a ...any) any {
						return w["a"]
					},
				),
			),
			xs: []T{
				{
					name:    "[1, [2, 3], 4]",
					in:      []any{1, []any{2, 3}, 4},
					matched: true,
					out:     3,
				},
				{
					name:    "[1, 2, 3, 4]",
					in:      []any{1, 2, 3, 4},
					matched: false,
					out:     nil,
				},
			},
		},

		{
			name: ",..",
			m: SyntaxRule(
				Case[any](
					L(1, 2, 3, ",..xs"),
					func(w W, a ...any) any {
						return w["xs"]
					},
				),
			),
			xs: []T{
				{
					name:    "[1, 2, 3]",
					in:      []any{1, 2, 3},
					matched: true,
					out:     []any{},
				},
				{
					name:    "[1, 2, 3, 4]",
					in:      []any{1, 2, 3, 4},
					matched: true,
					out:     []any{4},
				},
				{
					name:    "[1, 2, 3, 4, 5]",
					in:      []any{1, 2, 3, 4, 5},
					matched: true,
					out:     []any{4, 5},
				},
			},
		},
		{
			name: ",..",
			m: SyntaxRule(
				Case[any](
					L(1, 2, 3, ",.+xs"),
					func(w W, a ...any) any {
						return w["xs"]
					},
				),
			),
			xs: []T{
				{
					name:    "[1, 2, 3]",
					in:      []any{1, 2, 3},
					matched: false,
					out:     nil,
				},
				{
					name:    "[1, 2, 3, 4]",
					in:      []any{1, 2, 3, 4},
					matched: true,
					out:     []any{4},
				},
				{
					name:    "[1, 2, 3, 4, 5]",
					in:      []any{1, 2, 3, 4, 5},
					matched: true,
					out:     []any{4, 5},
				},
			},
		},
		{
			name: "Nested ,..",
			m: SyntaxRule(
				Case[any](
					L("let", L(S(",.+defs", L(S("lhs", P.Atom), ",rhs"))), ",body"),
					func(w W, ext ...any) any {
						return []any{w["defs"], w[",.+defs"], w["body"]}
					},
				),
			),
			xs: []T{
				{
					name: `
(let ([x 1]
	  [y 2])
	(+ 1 2))
`,
					in: []any{
						"let",
						[]any{
							[]any{"x", 1},
							[]any{"y", 2},
						},
						[]any{"+", 1, 2},
					},
					matched: true,
					out: []any{
						[]any{
							[]any{"x", 1},
							[]any{"y", 2},
						},
						[]W{
							{
								"lhs": "x",
								"rhs": 1,
							},
							{
								"lhs": "y",
								"rhs": 2,
							},
						},
						[]any{"+", 1, 2},
					},
				},
			},
		},
	}

	for _, c := range cases {
		for _, x := range c.xs {
			r, ok := c.m(x.in)
			Assert(ok == x.matched, c.name+"/"+x.name)
			Assert(reflect.DeepEqual(r, x.out), c.name+"/"+x.name)
		}
	}
}

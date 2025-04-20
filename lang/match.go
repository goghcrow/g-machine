package lang

import (
	"strings"
)

type (
	Ident = string
	Node  = any // int | string | []Node | []W
	W     = map[Ident]Node
)

type (
	Pattern      = any // int | string | *Slot | *RepeatPattern | BuiltPattern | []Pattern
	ListPattern  = []Pattern
	BuiltPattern = func(n Node, w W) bool
)

type (
	Matcher[T any] func(n Node, ext ...any) (T, bool)
	Slot           struct {
		ID         Ident
		SubPattern any // Pattern | BuiltPattern
	}
	SyntaxCase[T any] struct {
		XS []Pattern
		F  func(W, ...any) T
	}
)

func SyntaxRule[T any](cases ...SyntaxCase[T]) Matcher[T] {
	xss := SliceMap(cases, func(c SyntaxCase[T]) []BuiltPattern {
		return SliceMap(c.XS, buildPattern)
	})
	return func(n Node, ext ...any) (z T, _ bool) {
		for i, pts := range xss {
			for _, pt := range pts {
				w := W{}
				if pt(n, w) {
					return cases[i].F(w, ext...), true
				}
			}
		}
		return z, false
	}
}

func buildPattern(p any) BuiltPattern {
	switch p := p.(type) {
	case BuiltPattern:
		return p
	case *Slot:
		if p.SubPattern != nil {
			p.SubPattern = buildPattern(p.SubPattern)
		}
		return func(n Node, w W) bool {
			if p.SubPattern != nil && !p.SubPattern.(BuiltPattern)(n, w) {
				return false
			}
			w[p.ID] = n
			return true
		}
	case ListPattern:
		if sp, pre, rearID, rearN := endsDot(p); rearN >= 0 {
			xs := SliceMap(p[:len(p)-1], buildPattern)
			if sp != nil {
				sp = buildPattern(sp)
			}
			return func(n Node, w W) bool {
				ys, ok := n.([]Node)
				if !ok || len(ys) < len(xs)+rearN {
					return false
				}
				for i, x := range xs {
					if !x(ys[i], w) {
						return false
					}
				}

				var ws []W
				if sp != nil {
					ws = make([]W, len(ys[len(xs):]))
					for i, y := range ys[len(xs):] {
						ws[i] = W{}
						if !sp.(BuiltPattern)(y, ws[i]) {
							return false
						}
					}
				}
				w[rearID] = ys[len(xs):]
				w[pre+rearID] = ws // w[,..name] 指向嵌套结果
				return true
			}
		} else {
			xs := SliceMap(p, buildPattern)
			return func(n Node, w W) bool {
				ys, ok := n.([]Node)
				if !ok || len(ys) != len(xs) {
					return false
				}
				for i, x := range xs {
					if !x(ys[i], w) {
						return false
					}
				}
				return true
			}
		}
	case string:
		if len(p) > 0 && p[0] == ptnPrefix {
			return buildPattern(S(p[1:]))
		} else {
			return func(n Node, w W) bool {
				return n == p
			}
		}
	case int:
		return func(n Node, w W) bool {
			return n == p
		}
	default:
		panic("invalid pattern")
	}
}

func endsDot(p ListPattern) (ptn Pattern, dot, rearID string, n int) {
	if len(p) == 0 {
		return nil, dot, rearID, -1
	}
	switch s := p[len(p)-1].(type) {
	case string:
		pre, id, n := parseDot(s)
		return nil, pre, id, n
	case *Slot:
		pre, id, n := parseDot(s.ID)
		return s.SubPattern, pre, id, n
	}
	return nil, dot, rearID, -1
}

func parseDot(id string) (prefix, rearID string, n int) {
	if strings.HasPrefix(id, dot) {
		return dot, id[len(dot):], 0
	}
	if strings.HasPrefix(id, dotPlus) {
		return dotPlus, id[len(dotPlus):], 1
	}
	return dot, rearID, -1
}

const (
	ptnPrefix = ','
	dot       = ",.."
	dotPlus   = ",.+"
)

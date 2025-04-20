package lang

import "fmt"

type (
	Name  = string
	Arity = int
	Tag   = int
)

var Fmt = fmt.Sprintf

func Assert(cond bool, format string, a ...any) {
	if !cond {
		panic(fmt.Errorf(format, a...))
	}
}

func Iff[T any](cond bool, then, els T) T {
	if cond {
		return then
	}
	return els
}

func SliceMap[Slice ~[]A, A, B any](xs Slice, mapper func(A) B) []B {
	ys := make([]B, len(xs))
	for i, x := range xs {
		ys[i] = mapper(x)
	}
	return ys
}

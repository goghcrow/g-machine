package main

import "fmt"

type (
	Name   = string
	Offset = int
	Size   = int
	Arity  = int
	Tag    = int
	Addr   = int
)

var Fmt = fmt.Sprintf

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}
func Iff[T any](cond bool, then, els T) T {
	if cond {
		return then
	}
	return els
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func SliceMap[Slice ~[]A, A, B any](xs Slice, mapper func(A) B) []B {
	ys := make([]B, len(xs))
	for i, x := range xs {
		ys[i] = mapper(x)
	}
	return ys
}

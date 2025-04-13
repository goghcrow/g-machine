package main

type Stack[T any] struct {
	V []T
}

func (s *Stack[T]) stackSz() int             { return len(s.V) }
func (s *Stack[T]) empty() bool              { return len(s.V) == 0 }
func (s *Stack[T]) stackNth(offset Offset) T { return s.V[len(s.V)-1-offset] }
func (s *Stack[T]) peek() T                  { return s.V[len(s.V)-1] }
func (s *Stack[T]) push(v T)                 { s.V = append(s.V, v) }
func (s *Stack[T]) pushN(vs ...T)            { s.V = append(s.V, vs...) }
func (s *Stack[T]) pop() (v T) {
	Assert(len(s.V) > 0, "empty stack")
	topIdx := len(s.V) - 1
	v, s.V = s.V[topIdx], s.V[:topIdx]
	return
}
func (s *Stack[T]) popN(n int) []T {
	Assert(len(s.V) >= n, "invalid state")
	xs := make([]T, n)
	copy(xs, s.V[len(s.V)-n:])
	s.V = s.V[:len(s.V)-n]
	return xs
}
func (s *Stack[T]) drop(n int) {
	Assert(len(s.V) >= n, "invalid state")
	s.V = s.V[:len(s.V)-n]
}
func (s *Stack[T]) slide(n int) {
	top := len(s.V) - 1
	s.V[top-n] = s.V[top]
	s.V = s.V[:top]
}

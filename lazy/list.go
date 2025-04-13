package lazy

type Thunk[T any] func() T

type (
	List[T any] interface{ _list() }
	Nil[T any]  struct{}
	Cons[T any] struct {
		Car T
		Cdr Thunk[List[T]]
	}
)

func (Nil[T]) _list()  {}
func (Cons[T]) _list() {}

func ListFrom[T any](s []T) List[T] {
	var xs List[T] = Nil[T]{}
	for i := len(s) - 1; i >= 0; i-- {
		t := xs
		xs = Cons[T]{
			Car: s[i],
			Cdr: func() List[T] { return t },
		}
	}
	return xs
}
func ListMap[X, Y any](xs List[X], f func(X) Y) List[Y] {
	switch it := xs.(type) {
	case Nil[X]:
		return Nil[Y]{}
	case Cons[X]:
		return Cons[Y]{Car: f(it.Car), Cdr: func() List[Y] {
			return ListMap[X, Y](it.Cdr(), f)
		}}
	default:
		panic("unreachable")
	}
}

func ListTake[T any](xs List[T], n int) []T {
	var aux func(List[T], int, []T) []T
	aux = func(xs List[T], n int, acc []T) []T {
		if n <= 0 {
			return acc
		}
		switch it := xs.(type) {
		case Nil[T]:
			return acc
		case Cons[T]:
			return aux(it.Cdr(), n-1, append(acc, it.Car))
		default:
			panic("unreachable")
		}
	}
	return aux(xs, n, []T{})
}

func ListTo[T any](xs List[T]) []T {
	var aux func(List[T], []T) []T
	aux = func(xs List[T], acc []T) []T {
		switch it := xs.(type) {
		case Nil[T]:
			return acc
		case Cons[T]:
			return aux(it.Cdr(), append(acc, it.Car))
		default:
			panic("unreachable")
		}
	}
	return aux(xs, []T{})
}

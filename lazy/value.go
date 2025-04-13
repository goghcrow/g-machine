package lazy

type (
	LazyData[T any] interface{ _lazyD() }
	LazyRef[T any]  struct{ Data LazyData[T] }
)

type (
	Waiting[T any] Thunk[T]
	Done[T any]    struct{ Value T }
)

func (w Waiting[T]) _lazyD() {}
func (d Done[T]) _lazyD()    {}

func (r *LazyRef[T]) Extract() T {
	switch it := r.Data.(type) {
	case Waiting[T]:
		val := it()
		r.Data = Done[T]{val}
		return val
	case Done[T]:
		return it.Value
	default:
		panic("uhnreached")
	}
}

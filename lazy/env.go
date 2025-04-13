package lazy

type Env[K comparable, V any] func(name K) (V, bool)

func MkEnv[K comparable, V any]() Env[K, V] {
	return func(k K) (zero V, _ bool) {
		return zero, false
	}
}

func (e Env[K, V]) Extend(k K, v V) Env[K, V] {
	return func(k1 K) (V, bool) {
		if k == k1 {
			return v, true
		}
		return e(k1)
	}
}

func (e Env[K, V]) Lookup(k K) (V, bool) {
	return e(k)
}

package main

// 在编译超组合子时，需要维护一个环境，在编译过程中通过参数的名字找到参数在栈中的相对位置

type Env func(name Name) ( /*Stack*/ Offset, bool)

func MkEnv() Env {
	return func(name Name) (Offset, bool) { return -1, false }
}

func (e Env) Ext(name Name, offset Offset) Env {
	return func(lookup Name) (Offset, bool) {
		if name == lookup {
			return offset, true
		}
		return e(lookup)
	}
}

func (e Env) Lookup(name string) (int, bool) {
	return e(name)
}

func (e Env) Offset(n int) Env {
	Assert(n >= 0, "illegal state")
	if n == 0 {
		return e
	}
	return func(name Name) (Offset, bool) {
		if offset, ok := e(name); ok {
			return offset + n, ok
		}
		return -1, false
	}
}

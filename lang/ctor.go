package lang

import "fmt"

const (
	TagNil  = 1
	TagCons = 2
)

// 暂时不支持自定义数据结构
var builtinConstructors = []TCtor{
	Ctor("Nil", TagNil, 0),
	Ctor("Cons", TagCons, 2),
}

func ResolveCtor(nameOrTag any) TCtor {
	switch x := nameOrTag.(type) {
	case Tag:
		for _, ctor := range builtinConstructors {
			if ctor.Tag == x {
				return ctor
			}
		}
	case Name:
		for _, ctor := range builtinConstructors {
			if ctor.Name == x {
				return ctor
			}
		}
	}
	panic("undefined constructor " + fmt.Sprintf("%+v", nameOrTag))
}

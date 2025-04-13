package main

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

var preludeDefs = `
(define (i x) x)
(define (k x y) x)
(define (k1 x y) y)
(define (a f g x) (f x (g x)))
(define (compose f g x) (f (g x)))
(define (twice f) (compose f f))
`

//func() []TSC {
//	var id = SC("I", Vars("x"), Var("x"))                                                       // id x = x
//	var k = SC("K", Vars("x", "y"), Var("x"))                                                   // K x y = x
//	var k1 = SC("K1", Vars("x", "y"), Var("y"))                                                 // K1 x y = y
//	var s = SC("S", Vars("f", "g", "x"), App(App(Var("f"), Var("x")), App(Var("g"), Var("x")))) // S f g x = f x (g x)
//	var compose = SC("compose", Vars("f", "g", "x"), App(Var("f"), App(Var("g"), Var("x"))))    // compose f g x = f (g x)
//	var twice = SC("twice", Vars("f"), App(App(Var("compose"), Var("f")), Var("f")))            // twice f = compose f f
//	return []TSC{id, k, k1, s, compose, twice}
//}()

// ↓↓↓ eval 指令不一定总用得上, e.g. (add 3 (mul 4 5))
// `add`的两个参数在执行`Eval`之前就已经是WHNF, 这里的`Eval`指令是多余的

var compiledPrimitives = []ScDef{
	// 算术
	{"+", 2, []Instr{Push(1), IEval, Push(1), IEval, IAdd, Update(2), Pop(2), IUnwind}},
	{"-", 2, []Instr{Push(1), IEval, Push(1), IEval, ISub, Update(2), Pop(2), IUnwind}},
	{"*", 2, []Instr{Push(1), IEval, Push(1), IEval, IMul, Update(2), Pop(2), IUnwind}},
	{"/", 2, []Instr{Push(1), IEval, Push(1), IEval, IDiv, Update(2), Pop(2), IUnwind}},
	// 比较
	{"=", 2, []Instr{Push(1), IEval, Push(1), IEval, IEQ, Update(2), Pop(2), IUnwind}},
	{"!=", 2, []Instr{Push(1), IEval, Push(1), IEval, INE, Update(2), Pop(2), IUnwind}},
	{">=", 2, []Instr{Push(1), IEval, Push(1), IEval, IGE, Update(2), Pop(2), IUnwind}},
	{">", 2, []Instr{Push(1), IEval, Push(1), IEval, IGT, Update(2), Pop(2), IUnwind}},
	{"<=", 2, []Instr{Push(1), IEval, Push(1), IEval, ILE, Update(2), Pop(2), IUnwind}},
	{"<", 2, []Instr{Push(1), IEval, Push(1), IEval, ILT, Update(2), Pop(2), IUnwind}},
	// 杂项
	{"negate", 1, []Instr{Push(0), IEval, INeg, Update(1), Pop(1), IUnwind}},
	{"if", 3, []Instr{Push(0), IEval, Cond{
		Then: []Instr{Push(1)},
		Else: []Instr{Push(2)},
	}, Update(3), Pop(3), IUnwind}},
}

var builtinOps = func() map[Name]Code {
	m := map[Name]Code{}
	for _, prim := range compiledPrimitives {
		m[prim.Name] = prim.Code
	}
	return m
}()

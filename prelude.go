package g_machine

import (
	. "github.com/goghcrow/g_machine/lang"
)

var PreludeDefs = `
(define (I x) x)
(define (K x y) x)
(define (K1 x y) y)
(define (S f g x) (f x (g x)))
(define (compose f g x) (f (g x)))
(define (twice f) (compose f f))
`

// ↓↓↓ eval 指令不一定总用得上, e.g. (add 3 (mul 4 5))
// `add`的两个参数在执行`Eval`之前就已经是 WHNF, 这里的 `Eval` 指令是多余的

var CompiledPrimitives = []ScDef{
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
	for _, prim := range CompiledPrimitives {
		m[prim.Name] = prim.Code
	}
	return m
}()

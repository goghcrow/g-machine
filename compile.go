package g_machine

import (
	. "github.com/goghcrow/g_machine/lang"
)

// 在一个超组合子编译出的指令序列执行前，栈内一定已经存在这样一些地址：
//	最顶部的地址指向一个NGlobal节点(超组合子本身)
//	紧随其后的N个地址（N是该超组合子的参数数量）则指向一系列的App节点 - 正好对应到一个redex 的 spine
//	栈最底层的地址指向表达式最外层的 App节点，其余以此类推

// 通过调用 compileC 函数来生成对超组合子进行实例化的代码，并在后面加上三条指令。这三条指令各自的工作是：
//		Update(N)将堆中原本的redex更新为一个NInd节点，这个间接节点则指向刚刚实例化出来的超组合子
//		Pop(N)清理栈中已经无用的地址
//		Unwind寻找redex开始下一次规约

// CompileSC Compile Super Combinator
// 这里需要跟 Unwind(Global) 代码配合
func CompileSC(s TSC) ScDef {
	env := MkEnv()
	for i, arg := range s.Args {
		env = env.Ext(Name(arg), i)
	}
	arity := len(s.Args)
	var body Code
	if arity == 0 {
		// 指令Pop 0实际上什么也没做，故 arity == 0 时不生成
		body = append(compileC(s.Body, env), Update(arity), IUnwind)
	} else {
		body = append(compileC(s.Body, env), Update(arity), Pop(arity), IUnwind)
	}
	return ScDef{Name: s.Name, Arity: arity, Code: body}
}

func compileC(t Term, env Env) Code {
	switch t := t.(type) {
	case TVar:
		// 在编译超组合子的定义时使用比较粗糙的方式：
		// 一个变量如果不是参数，就当成其他超组合子（写错了会导致运行时错误）
		n, ok := env.Lookup(Name(t))
		if ok {
			//return []Instr{PushArg(n)}
			return []Instr{Push(n)}
		} else {
			return []Instr{PushGlobal(t)}
		}
	case TNum:
		return []Instr{PushInt(t)}
	case TCtor:
		ctor := ResolveCtor(t.Tag)
		return []Instr{Pack{Tag: ctor.Tag, Arity: ctor.Arity}}
	case TMatch:
		return compileE(t, env)
	case TApp:
		return MatchTerms[Code](t, []MatchAlts[Code]{
			{
				// Nil 不是 App, 这里只处理 (Cons ...)
				Apps(ResolveCtor(TagCons), TP.Slot("x"), TP.Slot("xs")),
				func(t Term, binds Binds) Code {
					ctor := ResolveCtor(TagCons)
					return append(
						append(
							compileC(binds["xs"], env),
							compileC(binds["x"], env.Offset(1))...,
						),
						Pack{Tag: ctor.Tag, Arity: ctor.Arity},
					)
				},
			},
			{
				Ptn: TP.Slot("_"),
				Fn: func(_ Term, binds Binds) Code {
					// 对于函数应用，先编译右侧表达式，然后将环境中所有参数对应的偏移量加一
					//（因为栈顶多出了一个地址指向实例化之后的右侧表达式），再编译左侧，最后加上MkApp指令
					return append(
						compileC(t.Arg, env),
						append(
							compileC(t.Fun, env.Offset(1)),
							IMkApp,
						)...,
					)
				},
			},
		})
	case TSC:
		//Pack
		panic("TODO")
	case TLet:
		return compileLet(compileC, t, env)
	default:
		panic("not support yet")
	}
}

func compileLet(compile func(Term, Env) Code, let TLet, env Env) (code Code) {
	if let.Rec {
		return compileLetRec(compile, let, env)
	}

	// 如果`(let (.....) e)`处于严格上下文中，那么`e`也处于严格上下文中
	// 但是前面的局部变量对应的表达式就不是，因为e不一定需要它们的结果
	// 所以编译 def 使用 compileC, 而编译 body 使用 compileE
	for _, def := range let.Defs {
		code = append(code, compileC(def.Val, env)...)
		// 更新偏移量并加入name所对应的本地变量的偏移量
		env = env.Offset(1).Ext(def.Var, 0)
	}
	return append(
		append(code, compile(let.Body, env)...),
		Slide(len(let.Defs)),
	)
}

func compileLetRec(compile func(Term, Env) Code, let TLet, env Env) (code Code) {
	//首先使用 Alloc(n) 申请N个地址
	//用 loop 表达式构建出完整的环境
	//编译 defs 中的本地变量，每编译完一个都用 Update 指令将结果更新到预分配的地址上
	//编译主表达式并用 Slide 指令清理现场

	for _, def := range let.Defs {
		env = env.Offset(1).Ext(def.Var, 0)
	}
	n := len(let.Defs)
	code = append(code, Alloc(n))
	for i := 0; i < n; i++ {
		code = append(code, compileC(let.Defs[i].Val, env)...)
		code = append(code, Update(n-1-i))
	}
	return append(
		append(code, compile(let.Body, env)...),
		Slide(n),
	)
}

func compileAlts(alts []TAlt, env Env) CaseJump {
	xs := make([]JumpCase, len(alts))
	for i, alt := range alts {
		env1 := env
		for i, v := range alt.Vars {
			env1 = env1.Ext(Name(v), i)
		}
		xs[i] = JumpCase{
			Tag:  alt.Ctor.Tag,
			Code: compileC(alt.Body, env1),
		}
	}
	return xs
}

// 一种可行的优化掉多余 eval 方法: 在编译表达式时注意其上下文
// e.g.，add需要它的参数被求值成WHNF，那么它的参数在编译时就处于严格(Strict)上下文中
// 通过这种方式，我们可以识别出一部分可以安全地按照严格求值进行编译的表达式(仅有一部分)
//
// 一个超组合子定义中的表达式处于严格上下文中
// 如果`(op e1 e2)`处于严格上下文中(此处`op`是一个primitive)，那么`e1`和`e2`也处于严格上下文中
// 如果`(let (.....) e)`处于严格上下文中，那么`e`也处于严格上下文中(但是前面的局部变量对应的表达式就不是，因为e不一定需要它们的结果)

// 严格求值上下文下的编译，它所生成的指令可以保证*栈顶地址指向的值一定是一个WHNF*。

func compileE(t Term, env Env) Code {
	switch t := t.(type) {
	case TNum:
		// 常数则直接push
		return []Instr{PushInt(t)}
	case TLet:
		// 编译一个严格上下文中的`let/letrec`表达式只需要用`compileE`编译其主表达式即可
		return compileLet(compileE, t, env)
	case TCtor:
		return compileC(t, env)
		// 这里应该不需要, 路由到 compileC
		//ctor := builtinConstructors[t.Name]
		//Assert(ctor != TCtor{}, "unsupported ctor "+t.Name)
		//return []Instr{
		//	Pack{Tag: ctor.Tag, Arity: ctor.Arity},
		//}
	case TMatch:
		return MatchTerms[Code](t, []MatchAlts[Code]{
			{
				Match(TP.Slot("e"), nil /*wildcard*/),
				func(t Term, binds Binds) Code {
					// 由于 case 表达式匹配的对象需要被求值到 WHNF，因此只能通过 compileE 函数来编译它
					return append(
						compileE(binds["e"], env),
						compileAlts(t.(TMatch).Alts, env),
					)
				},
			},
			{
				Ptn: TP.Slot("_"),
				Fn: func(t Term, binds Binds) Code {
					return compileC(t, env)
				},
			},
		})
	case TApp:
		return MatchTerms[Code](t, []MatchAlts[Code]{
			{
				Apps(Var("if"), TP.Slot("cond"), TP.Slot("then"), TP.Slot("else")),
				func(t Term, binds Binds) Code {
					return append(
						compileE(binds["cond"], env),
						Cond{
							Then: compileE(binds["then"], env),
							Else: compileE(binds["else"], env),
						},
					)
				},
			},
			{
				App(Var("negate"), TP.Slot("e")),
				func(t Term, binds Binds) Code {
					return append(compileE(binds["e"], env), INeg)
				},
			},
			{
				Apps(TP.Slot("op", isBuiltin), TP.Slot("lhs"), TP.Slot("rhs")),
				func(t Term, binds Binds) Code {
					code := builtinOps[Name(binds["op"].(TVar))]
					return append(
						append(
							compileE(binds["lhs"], env),
							compileE(binds["rhs"], env.Offset(1))...,
						),
						code...,
					)
				},
			},
			// 已经处理成了 TMatch, 不是 TAPPs
			//{
			//	Match(TP.Slot("e"), nil /*wildcard*/),
			//	func(t Term, binds Binds) Code {
			//		// 由于 case 表达式匹配的对象需要被求值到 WHNF，因此只能通过 compileE 函数来编译它
			//		return append(
			//			compileE(binds["e"], env),
			//			compileAlts(t.(TMatch).Alts, env),
			//		)
			//	},
			//},
			// 这里应该不需要, 路由到 compileC
			//{
			//	// Nil 不是 App, 这里只处理 (cons ...)
			//	Apps(builtinConstructors["Cons"], TP.Slot("x"), TP.Slot("xs")),
			//	func(t Term, binds Binds) Code {
			//		consCtor := builtinConstructors["Cons"]
			//		return append(
			//			append(
			//				compileC(binds["xs"], env),
			//				compileC(binds["x"], env.Offset(1))...,
			//			),
			//			Pack{Tag: consCtor.Tag, Arity: consCtor.Arity},
			//		)
			//	},
			//},
			{
				Ptn: TP.Slot("_"),
				Fn: func(t Term, binds Binds) Code {
					return compileC(t, env)
				},
			},
		})
	case TSC:
		//Pack
		panic("TODO")
	default:
		// 默认分支，我们仅仅在`compileC`的结果后面加一条`Eval`指令
		return append(compileC(t, env), IEval)
	}
}

func isBuiltin(t Term) bool {
	if v, ok := t.(TVar); ok {
		return builtinOps[Name(v)] != nil
	}
	return false
}

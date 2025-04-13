package main

import "fmt"

type Term interface {
	isTerm()
	fmt.Stringer
}

// todo term 应该是 TLam, 然后 lambda lift 成 super combinator
// 这里直接 define super combinator

type (
	TVar  Name // Identity
	TNum  int
	TCtor struct {
		Name
		Tag
		Arity
	}
	TApp struct{ Fun, Arg Term }
	TSC  struct { // Super Combinator
		Name string
		Args []TVar // Params
		Body Term
	}

	TCase struct {
		Expr Term
		Alts []TAlt
	}
	TLet struct {
		Rec  bool
		Defs []TDef
		Body Term
	}

	TDef struct {
		Var Name
		Val Term
	}
	TAlt struct {
		Ctor TCtor
		Vars []TVar
		Body Term
	}
)

func (TVar) isTerm()  {}
func (TNum) isTerm()  {}
func (TCtor) isTerm() {}
func (TApp) isTerm()  {}
func (TSC) isTerm()   {}
func (TLet) isTerm()  {}
func (TCase) isTerm() {}

func Var(name Name) TVar { return TVar(name) }
func Vars(names ...Name) (xs []TVar) {
	for _, name := range names {
		xs = append(xs, Var(name))
	}
	return
}
func Num(n int) TNum { return TNum(n) }
func Ctor(name Name, tag Tag, arity Arity) TCtor {
	return TCtor{
		Name:  name,
		Tag:   tag,
		Arity: arity,
	}
}
func App(fun, arg Term) TApp {
	return TApp{Fun: fun, Arg: arg}
}
func Apps(fun Term, arg Term, rest ...Term) TApp {
	app := TApp{Fun: fun, Arg: arg}
	for _, x := range rest {
		app = TApp{Fun: app, Arg: x}
	}
	return app
}
func Def(name Name, val Term) TDef { return TDef{Var: name, Val: val} }
func Let(isRec bool, defs []TDef, body Term) TLet {
	return TLet{Rec: isRec, Defs: defs, Body: body}
}
func Case(expr Term, alts []TAlt) TCase {
	return TCase{Expr: expr, Alts: alts}
}
func Alt(ctor TCtor, vars []TVar, body Term) TAlt {
	return TAlt{Ctor: ctor, Vars: vars, Body: body}
}
func SC(name Name, args []TVar, body Term) TSC {
	return TSC{Name: name, Args: args, Body: body}
}

func (t TVar) String() string  { return Name(t) }
func (t TNum) String() string  { return Fmt("%d", t) }
func (t TCtor) String() string { return t.Name }
func (t TApp) String() string  { return Fmt("(%s %s)", t.Fun, t.Arg) }
func (t TSC) String() string   { return Fmt("(define (%s %s) %s)", t.Name, t.Args, t.Body) }
func (t TLet) String() string  { return Fmt("(let%s (%s) %s)", Iff(t.Rec, "rec", ""), t.Defs, t.Body) }
func (t TCase) String() string { return Fmt("(case %s %s)", t.Expr, t.Alts) }
func (t TDef) String() string  { return Fmt("[%s %s]", t.Var, t.Val) }
func (t TAlt) String() string  { return Fmt("[(%s %s) %s]", t.Ctor, t.Vars, t.Body) }

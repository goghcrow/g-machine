package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token = string

var enclose = map[Token]Token{
	"(": ")",
	"[": "]",
}

func isStarting(tok Token) bool {
	for starting := range enclose {
		if tok == starting {
			return true
		}
	}
	return false
}
func isClosing(tok Token) bool {
	for _, closing := range enclose {
		if tok == closing {
			return true
		}
	}
	return false
}

func Parse(src string) (pgrm []Term) {
	toks := tokenize(src)
	var ast0 any
	for len(toks) > 0 {
		ast0, toks = parse0(toks)
		pgrm = append(pgrm, parse1(ast0))
	}
	return
}

func tokenize(src string) []Token {
	for starting, closing := range enclose {
		src = strings.ReplaceAll(src, starting, " "+starting+" ")
		src = strings.ReplaceAll(src, closing, " "+closing+" ")
	}
	return strings.Fields(src)
}

func parse0(toks []Token) (any, []Token) {
	Assert(len(toks) > 0, "unexpected EOF")

	car, cdr := toks[0], toks[1:]
	Assert(!isClosing(car), "unexpected "+car)

	if isStarting(car) {
		Assert(len(cdr) > 0, "unexpected EOF")
		//goland:noinspection GoPreferNilSlice
		lst := []any{} // 确保不是 nil
		for !isClosing(cdr[0]) {
			var exp any
			exp, cdr = parse0(cdr)
			lst = append(lst, exp)
		}
		return lst, cdr[1:]
	}

	// TNum | TCtor | TVar
	i, err := strconv.ParseInt(car, 10, 64)
	if err == nil {
		return Num(int(i)), cdr
	}
	// 用大小写区分 Constructor 和 Variable
	r, _ := utf8.DecodeRuneInString(car)
	if unicode.IsUpper(r) {
		return ResolveCtor(car), cdr
	} else {
		return Var(car), cdr
	}
}

func parse1(ast0 any) Term {
	switch xs := ast0.(type) {
	case []any:
		Assert(len(xs) > 0, "unexpected ()")
		switch xs[0] {
		case Var("define"): // TSC
			return parseDefine(xs)
		case Var("case"): // TCase
			return parseCase(xs)
		case Var("let"), Var("letrec"): // TLet
			return parseLet(xs)
		default: // TApp
			Assert(len(xs) > 1, "unexpected zero arg app")
			app := TApp{Fun: parse1(xs[0]), Arg: parse1(xs[1])}
			for _, x := range xs[2:] {
				app = TApp{Fun: app, Arg: parse1(x)}
			}
			return app
		}
	case TVar, TCtor, TNum:
		return xs.(Term)
	default:
		Assert(false, "unexpected "+fmt.Sprintf("%+v", xs))
		return nil
	}
}

// (define (name x...) exp)
// (define name exp)
func parseDefine(xs []any) TSC {
	Assert(len(xs) == 3, "invalid define form")
	switch x := xs[1].(type) {
	case TVar:
		return SC(Name(x), []TVar{}, parse1(xs[2]))
	case []any:
		id, ok := x[0].(TVar)
		Assert(ok && len(x) > 0, "invalid define name")
		return SC(Name(id), parseVars(x[1:]), parse1(xs[2]))
	default:
		Assert(false, "invalid define form")
		return TSC{}
	}
}

// (case a [(Nil) Nil] [(Cons x xs) exp] ...)
func parseCase(xs []any) TCase {
	Assert(len(xs) > 2, "invalid case form")

	expr := parse1(xs[1])

	var alts []TAlt
	for _, c := range xs[2:] {
		// [(Constructor var...) exp]
		ys, ok := c.([]any)
		Assert(ok && len(ys) == 2, "invalid case alt")

		ptn, ok := ys[0].([]any)
		Assert(ok && len(ptn) > 0, "invalid case alt")

		ctor, ok := ptn[0].(TCtor)
		Assert(ok, "invalid case alt constructor")

		alts = append(alts, Alt(ctor, parseVars(ptn[1:]), parse1(ys[1])))
	}

	return Case(expr, alts)
}

// (let ([a exp] [b exp] ...) body)
// (letrec  ([a exp] [b exp] ...) body)
func parseLet(xs []any) TLet {
	Assert(len(xs) == 3, "invalid let form")

	defs, ok := xs[1].([]any)
	Assert(ok, "invalid let def")
	tDefs := make([]TDef, len(defs))
	for i, def := range defs {
		def, ok := def.([]any)
		Assert(ok && len(def) == 2, "invalid let def")
		name, ok := def[0].(string)
		Assert(ok, "invalid let def var")
		tDefs[i] = Def(name, parse1(def[1]))
	}

	return Let(xs[0] == "letrec", tDefs, parse1(xs[2]))
}

func parseVars(xs []any) []TVar {
	ok := false
	vars := make([]TVar, len(xs))
	for i, v := range xs {
		vars[i], ok = v.(TVar)
		Assert(ok, "invalid vars")
	}
	return vars
}

package lang

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	Token = string
	Pgrm  = []TSC
)

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

func Parse(src string) (pgrm Pgrm) {
	initSyntax()
	toks := tokenize(src)
	var ast0 any
	for len(toks) > 0 {
		ast0, toks = parse0(toks)
		term := parse1(ast0)
		tsc, ok := term.(TSC)
		Assert(ok, "want sc but %T", term)
		pgrm = append(pgrm, tsc)
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
		lst := []Node{} // 确保不是 nil
		for !isClosing(cdr[0]) {
			var exp any
			exp, cdr = parse0(cdr)
			lst = append(lst, exp)
			Assert(len(cdr) > 0, "unexpected EOF")
		}
		return lst, cdr[1:]
	}

	// TNum | TCtor | TVar
	i, err := strconv.ParseInt(car, 10, 64)
	if err == nil {
		return int(i), cdr
	}
	return car, cdr
}

func parse1(ast0 any) Term {
	switch xs := ast0.(type) {
	case []Node:
		Assert(len(xs) > 0, "syntax err: empty")
		switch xs[0] {
		case "define": // TSC
			tDefine, ok := matchDefine(xs)
			Assert(ok, "syntax error: define")
			return tDefine
		case "case": // TMatch
			tCase, ok := matchCase(xs)
			Assert(ok, "syntax error: case")
			return tCase
		case "let", "letrec": // TLet
			tLet, ok := matchLet(xs)
			Assert(ok, "syntax error: let")
			return tLet
		default: // TApp
			app, ok := matchApp(xs)
			Assert(ok, "syntax error: app")
			return app
		}
	case int:
		return TNum(xs)
	case string:
		return TVar(xs)
	default:
		Assert(false, "unexpected "+fmt.Sprintf("%+v", xs))
		return nil
	}
}

var (
	matchDefine Matcher[Term] = nil
	matchCase   Matcher[Term] = nil
	matchLet    Matcher[Term] = nil
	matchApp    Matcher[Term] = nil
)

func initSyntax() {
	if matchDefine == nil { // once?
		matchDefine = defineSyntax()
		matchCase = caseSyntax()
		matchLet = letSyntax()
		matchApp = appSyntax()
	}
}

func defineSyntax() Matcher[Term] {
	// (define (name x...) exp)
	// (define name exp)
	return SyntaxRule[Term](
		Case[Term](
			L("define", S("name", P.Atom), ",body"),
			func(w W, ext ...any) Term {
				return SC(w["name"].(Name), []TVar{}, parse1(w["body"]))
			},
		),
		Case[Term](
			L("define", L(S("name", P.Atom), S(",..args", P.Atom)), ",body"),
			func(w W, ext ...any) Term {
				return SC(
					w["name"].(Name),
					SliceMap(w["args"].([]Node), func(a Node) TVar { return TVar(a.(Name)) }),
					parse1(w["body"]),
				)
			},
		),
	)
}

func caseSyntax() Matcher[Term] {
	// (case a [(Nil) Nil] [(Cons x xs) exp] ...)
	return SyntaxRule[Term](
		Case[Term](
			L("case", ",a",
				S(",.+alts",
					L(
						L(S("ctor", P.Ctor), S(",..args", P.Atom)),
						",exp",
					),
				)),
			func(w W, ext ...any) Term {
				alts := w[",.+alts"].([]W)
				tAlts := make([]TAlt, len(alts))
				for i, alt := range alts {
					tAlts[i] = Alt(
						ResolveCtor(alt["ctor"]),
						SliceMap(alt["args"].([]Node), func(a Node) TVar { return TVar(a.(Name)) }),
						parse1(alt["exp"]),
					)
				}
				return Match(parse1(w["a"]), tAlts)
			},
		),
	)
}

func letSyntax() Matcher[Term] {
	// (let ([a exp] [b exp] ...) body)
	// (letrec  ([a exp] [b exp] ...) body)
	return SyntaxRule[Term](
		Case[Term](
			L(",let",
				L(
					S(",.+defs", L(S("lhs", P.Atom), ",rhs")),
				),
				",body"),
			func(w W, ext ...any) Term {
				defs := w[",.+defs"].([]W)
				tDefs := make([]TDef, len(defs))
				for i, def := range defs {
					tDefs[i] = Def(def["lhs"].(Name), parse1(def["rhs"]))
				}
				return Let(w["let"] == "letrec", tDefs, parse1(w["body"]))
			},
		),
	)
}

func appSyntax() Matcher[Term] {
	return SyntaxRule[Term](
		Case[Term](
			L(S("ctor", P.Ctor), ",..args"),
			func(w W, a ...any) Term {
				ctor := ResolveCtor(w["ctor"].(Name))
				args := SliceMap(w["args"].([]Node), func(a Node) Term { return parse1(a) })
				return Apps(ctor, args...)
			},
		),
		Case[Term](
			L(S("sc"), ",..args"),
			func(w W, a ...any) Term {
				f := parse1(w["sc"])
				args := SliceMap(w["args"].([]Node), func(a Node) Term { return parse1(a) })
				return Apps(f, args...)
			},
		),
	)
}

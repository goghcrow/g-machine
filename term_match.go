package main

type Pattern = Term

type (
	PVar      struct{ Name }
	PGuardVar struct {
		Name
		Guard func(Term) bool
	}
	PFun func(t Term, binds Binds) bool
)

type Binds = map[Name]Term

func (PVar) isTerm()      {}
func (PGuardVar) isTerm() {}
func (PFun) isTerm()      {}

func (p PVar) String() string      { return Fmt("PVar(%s)", p.Name) }
func (p PGuardVar) String() string { return Fmt("PGuardVar(%s)", p.Name) }
func (p PFun) String() string      { return Fmt("PFun") }

type patternT int

const P patternT = iota

func (patternT) Var(name Name) PVar { return PVar{name} }
func (patternT) Where(f PFun) PFun  { return f }
func (v PVar) Where(f func(t Term) bool) PGuardVar {
	return PGuardVar{v.Name, f}
}

type MatchAlts[R any] struct {
	Ptn Pattern
	Fn  func(Term, Binds) R
}

func MatchTerms[R any](t Term, alts []MatchAlts[R]) R {
	for _, alt := range alts {
		binds := Binds{}
		if matchTerm(t, alt.Ptn, binds) {
			return alt.Fn(t, binds)
		}
	}
	panic("illegal state")
}

func matchTerms(ts []Term, ptns []Pattern, binds Binds) bool {
	if ptns == nil /*wildcard*/ {
		return true
	}
	if len(ts) != len(ptns) {
		return false
	}
	for i, t := range ts {
		if !matchTerm(t, ptns[i], binds) {
			return false
		}
	}
	return true
}

func matchDefs(defs, ptn []TDef, binds Binds) bool {
	if ptn /*wildcard*/ == nil {
		return true
	}
	if len(defs) != len(ptn) {
		return false
	}
	for i, tDef := range defs {
		pDef := ptn[i]
		if pDef.Var != "" /*wildcard*/ && tDef.Var != pDef.Var {
			return false
		}
		if !matchTerm(tDef.Val, pDef.Val, binds) {
			return false
		}
	}
	return true
}

func matchAlts(alts, ptn []TAlt, binds Binds) bool {
	if ptn /*wildcard*/ == nil {
		return true
	}
	if len(alts) != len(ptn) {
		return false
	}
	for i, tc := range alts {
		pc := ptn[i]
		if !matchTerm(tc.Ctor, pc.Ctor, binds) {
			return false
		}
		if pc.Vars /*wildcard*/ != nil {
			if len(tc.Vars) != len(pc.Vars) {
				return false
			}
			for i, tv := range tc.Vars {
				pv := pc.Vars[i]
				if Name(tv) != "" && tv != pv {
					return false
				}
			}
		}
		if !matchTerm(tc.Body, pc.Body, binds) {
			return false
		}
	}
	return true
}

func matchTerm(t Term, ptn Pattern, binds Binds) bool {
	do := func(t Term, ptn Term) bool {
		return matchTerm(t, ptn, binds)
	}

	switch p := ptn.(type) {
	case PVar:
		binds[p.Name] = t
		return true
	case PGuardVar:
		if p.Guard(t) {
			binds[p.Name] = t
			return true
		}
		return false
	case PFun:
		return p(t, binds)
	}

	switch t := t.(type) {
	case TVar:
		switch p := ptn.(type) {
		case TVar:
			return t == p
		default:
			return false
		}
	case TNum:
		switch p := ptn.(type) {
		case TNum:
			return t == p
		default:
			return false
		}
	case TCtor:
		switch p := ptn.(type) {
		case TCtor:
			return (p.Name == "" /*wildcard*/ || t.Name == p.Name) &&
				(p.Tag == 0 /*wildcard*/ || t.Tag == p.Tag) &&
				(p.Arity == -1 /*wildcard*/ || t.Arity == p.Arity)
		default:
			return false
		}
	case TApp:
		switch p := ptn.(type) {
		case TApp:
			return do(t.Fun, p.Fun) && do(t.Arg, p.Arg)
		default:
			return false
		}
	case TSC:
		switch p := ptn.(type) {
		case TSC:
			return (p.Name == "" /*wildcard*/ || t.Name == p.Name) &&
				matchTerms(
					SliceMap(t.Args, func(it TVar) Term { return it }),
					SliceMap(p.Args, func(it TVar) Term { return it }),
					binds,
				) && do(t.Body, p.Body)
		default:
			return false
		}
	case TLet:
		switch p := ptn.(type) {
		case TLet:
			// isRec 不支持通配符, 用 PFun 匹配
			return t.Rec == p.Rec && matchDefs(t.Defs, p.Defs, binds) && do(t.Body, p.Body)
		default:
			return false
		}
	case TCase:
		switch p := ptn.(type) {
		case TCase:
			return do(t.Expr, p.Expr) && matchAlts(t.Alts, p.Alts, binds)
		default:
			return false
		}
	default:
		panic("illegal state")
	}
}

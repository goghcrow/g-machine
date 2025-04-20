package lang

func S(id Ident, optP ...Pattern) *Slot {
	switch len(optP) {
	case 0:
		return &Slot{ID: id}
	case 1:
		return &Slot{ID: id, SubPattern: optP[0]}
	default:
		panic("invalid args")
	}
}

func L(xs ...Pattern) ListPattern {
	return xs
}

const P pattern = iota

type pattern int

func (pattern) Atom(n Node, _ W) bool {
	_, ok := n.(string)
	return ok
}
func (pattern) Ctor(n Node, w W) bool {
	//r, _ := utf8.DecodeRuneInString(car)
	//unicode.IsUpper(r)
	return P.Atom(n, w) && len(n.(string)) > 0 && n.(string)[0] >= 'A' && n.(string)[0] <= 'Z'
}
func (pattern) Var(n Node, _ W) bool {
	switch n := n.(type) {
	case string:
		return true
	case ListPattern:
		return len(n) > 0 && n[0] == ".id"
	default:
		return false
	}
}

func (pattern) Any(n Node, _ W) bool { return true }

func (pattern) Prim(n Node, w W) bool {
	return P.Atom(n, w) && len(n.(string)) > 0 && n.(string)[0] == '.'
}

func (pattern) Trival(n Node, _ W) bool {
	switch n := n.(type) {
	case string:
		return true
	case ListPattern:
		return len(n) > 0 && (n[0] == ".quote" || n[0] == ".id" || n[0] == ".nil")
	default:
		return false
	}
}

func Case[T any](p Pattern, f func(W, ...any) T) SyntaxCase[T] {
	return SyntaxCase[T]{
		XS: []Pattern{p},
		F:  f,
	}
}

func Cases[T any](ps []Pattern, f func(W, ...any) T) SyntaxCase[T] {
	return SyntaxCase[T]{
		XS: ps,
		F:  f,
	}
}

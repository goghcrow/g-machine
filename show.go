package g_machine

import (
	"strings"

	. "github.com/goghcrow/g_machine/lang"
	"github.com/goghcrow/go-ansi"
)

func (s *State) showAddr(addr Addr) string {
	return s.showNode(s.Read(addr))
}

func (s *State) showCode() string {
	switch len(s.Code) {
	case 0:
		return ""
	case 1:
		car := s.Code[0].String()
		return ansi.Blue.Bold().Text(car).S()
	default:
		car, cdr := s.Code[0].String(), strings.Join(SliceMap(s.Code[1:], Instr.String), ", ")
		return ansi.Blue.Bold().Text(car).S() + ansi.Yellow.Text(", "+cdr).S()
	}
}

func (s *State) showStack() string {
	buf := ansi.Purple.Bold().Text("Stack").S()
	for _, addr := range s.Stack.V {
		buf += "\n\t" + s.showAddr(addr)
	}
	return buf
}

func (s *State) showNode(n GNode) string {
	switch n := n.(type) {
	case NNum:
		return Fmt("Num(%d)", n)
	case NApp:
		return Fmt("App(%s, %s)", s.showAddr(n.Fun), s.showAddr(n.Arg))
	case NInd:
		return Fmt("Ind(%s)", s.showAddr(Addr(n)))
	case NGlobal:
		return Fmt("Global(%s)", n.Name)
	case NCtor:
		return Fmt("Constructor(%d, %s)", n.Tag,
			"["+strings.Join(SliceMap(n.Args, s.showAddr), ",")+"]")
	default:
		panic("Unknown node type")
	}
}

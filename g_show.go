package main

import (
	"strings"
)

func (s *State) showAddr(addr Addr) string {
	return s.showNode(s.Read(addr))
}

func (s *State) showPC() string {
	if len(s.Code) == 0 {
		return ""
	}
	in := s.Code[0]
	if _, ok := in.(Unwind); ok {
		return Fmt("%s(%s)", in, s.showAddr(s.peek()))
	} else {
		return Fmt("%s", in)
	}
}
func (s *State) showCode() string {
	return "[" + strings.Join(SliceMap(s.Code, Instr.String), ",") + "]"
}

func (s *State) showStack() string {
	buf := "Stack"
	for _, addr := range s.Stack.V {
		buf += "\n\t" + s.showAddr(addr)
	}
	return buf
}

func (s *State) showNode(n Node) string {
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

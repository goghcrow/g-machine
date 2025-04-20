package g_machine

import (
	"fmt"
	"strings"

	. "github.com/goghcrow/g_machine/lang"
)

type Code = []Instr

// Instr Instruction
type Instr interface {
	isInstr() // sealed
	fmt.Stringer
}

var (
	IPrint  = Print{}
	IEval   = Eval{}
	IUnwind = Unwind{}
	IMkApp  = MkApp{}
	IAdd    = Add{}
	ISub    = Sub{}
	IMul    = Mul{}
	IDiv    = Div{}
	INeg    = Neg{}
	IEQ     = Eq{}
	INE     = Ne{}
	ILT     = Lt{}
	ILE     = Le{}
	IGT     = Gt{}
	IGE     = Ge{}
)

type (
	JumpCase = struct {
		Tag
		Code
	}
	JumpTable   = []JumpCase
	Constructor = struct{ Tag, Arity int }
)

type (
	Print      struct{}
	Eval       struct{}
	Unwind     struct{}
	PushGlobal Name
	PushInt    int
	Push       Offset
	PushArg    Offset // todo delete
	MkApp      struct{}
	Update     Offset
	Pop        Size
	Slide      Size
	Alloc      Size
	Split      Arity
	Pack       Constructor
	CaseJump   JumpTable

	Add struct{}
	Sub struct{}
	Mul struct{}
	Div struct{}
	Neg struct{}
	Eq  struct{}
	Ne  struct{}
	Lt  struct{}
	Le  struct{}
	Gt  struct{}
	Ge  struct{}

	Cond struct{ Then, Else Code }
)

func (Print) isInstr()      {}
func (Eval) isInstr()       {}
func (Unwind) isInstr()     {}
func (PushGlobal) isInstr() {}
func (PushInt) isInstr()    {}
func (Push) isInstr()       {}
func (PushArg) isInstr()    {}
func (MkApp) isInstr()      {}
func (Update) isInstr()     {}
func (Pop) isInstr()        {}
func (Slide) isInstr()      {}
func (Alloc) isInstr()      {}
func (Split) isInstr()      {}
func (Pack) isInstr()       {}
func (CaseJump) isInstr()   {}
func (Add) isInstr()        {}
func (Sub) isInstr()        {}
func (Mul) isInstr()        {}
func (Div) isInstr()        {}
func (Neg) isInstr()        {}
func (Eq) isInstr()         {}
func (Ne) isInstr()         {}
func (Lt) isInstr()         {}
func (Le) isInstr()         {}
func (Gt) isInstr()         {}
func (Ge) isInstr()         {}
func (Cond) isInstr()       {}

func (Print) String() string        { return "Print" }
func (Eval) String() string         { return "Eval" }
func (Unwind) String() string       { return "Unwind" }
func (i PushGlobal) String() string { return Fmt("PushGlobal(%s)", string(i)) }
func (i PushInt) String() string    { return Fmt("PushInt(%d)", i) }
func (i Push) String() string       { return Fmt("Push(%d)", i) }
func (i PushArg) String() string    { return Fmt("PushArg(%d)", i) }
func (MkApp) String() string        { return "MkApp" }
func (i Update) String() string     { return Fmt("Update(%d)", i) }
func (i Pop) String() string        { return Fmt("Pop(%d)", i) }
func (i Slide) String() string      { return Fmt("Slide(%d)", i) }
func (i Alloc) String() string      { return Fmt("Alloc(%d)", i) }
func (i Split) String() string      { return Fmt("Split(%d)", i) }
func (i Pack) String() string       { return Fmt("Pack(%d, %d)", i.Tag, i.Arity) }
func (Add) String() string          { return "Add" }
func (Sub) String() string          { return "Sub" }
func (Mul) String() string          { return "Mul" }
func (Div) String() string          { return "Div" }
func (Neg) String() string          { return "Neg" }
func (Eq) String() string           { return "Eq" }
func (Ne) String() string           { return "Ne" }
func (Lt) String() string           { return "Lt" }
func (Le) String() string           { return "Le" }
func (Gt) String() string           { return "Gt" }
func (Ge) String() string           { return "Ge" }
func (i Cond) String() string       { return Fmt("Cond(%s, %s)", showCode(i.Then), showCode(i.Else)) }
func (i CaseJump) String() string {
	return Fmt("CaseJump(%v)", SliceMap(i, func(a JumpCase) Tag { return a.Tag }))
}

func showCode(xs Code) string {
	buf := make([]string, len(xs))
	for i, x := range xs {
		buf[i] = x.String()
	}
	return "[" + strings.Join(buf, ", ") + "]"
}

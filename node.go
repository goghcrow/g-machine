package g_machine

import (
	. "github.com/goghcrow/g_machine/lang"
)

type ScDef struct {
	Name
	Arity
	Code
}

// GNode 代表 Graph Node
type GNode interface{ isNode() } // sealed

type (
	NNum    int
	NApp    struct{ Fun, Arg Addr }
	NInd    Addr  // Indirection节点, 实现惰性求值的关键一环
	NGlobal ScDef // 存放超组合子的参数数量和对应指令序列
	NCtor   struct {
		Tag
		Args []Addr
	}
)

func (NNum) isNode()    {}
func (NInd) isNode()    {}
func (NApp) isNode()    {}
func (NGlobal) isNode() {}
func (NCtor) isNode()   {}

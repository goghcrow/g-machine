package g_machine

import (
	"fmt"

	. "github.com/goghcrow/g_machine/lang"
)

type (
	Globals = map[Name]Addr
	Dump    struct {
		Code
		*Stack[Addr]
	}
)

func MkStack(xs ...Addr) *Stack[Addr] {
	return &Stack[Addr]{V: xs}
}

type State struct {
	Dump *Stack[Dump]
	Code []Instr
	*Stack[Addr]
	*Heap
	Globals Globals
	Output  []GNode
	// Stats   int
}

func (s *State) output(n GNode) {
	s.Output = append(s.Output, n)
}

func (s *State) lookUp(name string) Addr {
	addr, ok := s.Globals[name]
	if !ok {
		panic("undefined " + name)
	}
	return addr
}

func (s *State) allocNodes(n int) {
	for i := 0; i < n; i++ {
		s.push(s.Alloc(NInd(-1)))
	}
}

func (s *State) pushInstr(is ...Instr) { s.Code = append(is, s.Code...) }
func (s *State) popInstr() (in Instr) {
	in, s.Code = s.Code[0], s.Code[1:]
	return in
}

func (s *State) step() bool {
	if len(s.Code) == 0 {
		return false
	}

	switch i := s.popInstr().(type) {
	case Print:
		switch n := s.Read(s.pop()).(type) {
		case NNum:
			s.output(n)
		case NCtor:
			switch n.Tag {
			case TagCons:
				// 需要强制对addr1和addr2进行求值，故先执行Eval指令
				s.pushInstr(IEval, IPrint, IEval, IPrint)

				// todo todo todo
				// NCtor(1, Cons(addr1, Cons(addr2, Nil)))
				addr1 := n.Args[0]
				addr2 := s.Read(n.Args[1]).(NCtor).Args[0]
				s.push(addr2)
				s.push(addr1)
			case TagNil:
				s.output(n)
			default:
				panic("illegal state")
			}
		default:
			panic("illegal state")
		}

	case Eval:
		// 首先弹出栈顶地址
		// 然后保存(dump)当前还没执行的指令序列和栈
		// 清空当前栈并放入之前保存的地址
		// 清空当前指令序列，放入指令 Unwind
		addr := s.pop()
		s.Dump.push(Dump{Code: s.Code, Stack: s.Stack})
		s.Stack = MkStack(addr)
		s.Code = []Instr{IUnwind}

	case Unwind:
		s.stepUnwind()
	case PushGlobal:
		s.push(s.lookUp(Name(i)))
	case PushInt:
		s.push(s.Alloc(NNum(i)))
	case Push:
		// 将第offset + 1个地址复制到栈顶
		//    Push(n) a0 : . . . : an : s
		// => an : a0 : . . . : an : s
		s.push(s.stackNth(Offset(i)))
	case PushArg:
		// 栈地址布局：第一个地址应该指向超组合子节点，紧随其后的n个地址则指向N个NApp节点
		s.push(s.Read(s.stackNth(Offset(i) + 1 /*skip SC*/)).(NApp).Arg)
	case MkApp:
		app := s.Alloc(NApp{Fun: s.pop(), Arg: s.pop()})
		s.push(app)
	case Update:
		// 假设栈内第一个地址指向当前redex求值结果
		addr := s.pop()
		// 跳过紧随其后的超组合子节点地址
		// 把第N个NApp节点替换为一个指向求值结果的间接节点
		// 如果当前redex是CAF，那就直接把它在堆上的NGlobal节点替换掉
		dst := s.stackNth(Offset(i))
		s.Write(dst, NInd(addr))
	case Pop:
		s.drop(Size(i))
	case Slide:
		s.slide(Size(i))
	case Alloc:
		s.allocNodes(Size(i))
	case Cond:
		switch s.Read(s.pop()).(NNum) {
		case 1:
			s.pushInstr(i.Then...)
		case 0:
			s.pushInstr(i.Else...)
		default:
			panic("illegal state")
		}
	case Split:
		ctor := s.Read(s.pop()).(NCtor)
		Assert(Arity(i) == len(ctor.Args), "illegal state")
		s.pushN(s.Read(s.pop()).(NCtor).Args...)
	case Pack:
		s.push(s.Alloc(NCtor{i.Tag, s.popN(i.Arity)}))
	case CaseJump:
		addr := s.pop()
		t := s.Read(addr).(NCtor).Tag
		for _, jmp := range i {
			if jmp.Tag == t {
				s.pushInstr(jmp.Code...)
				s.push(addr)
				return true
			}
		}
		panic("illegal state")
	case Add:
		s.liftArithmetic2(func(x, y NNum) GNode { return x + y })
	case Sub:
		s.liftArithmetic2(func(x, y NNum) GNode { return x - y })
	case Mul:
		s.liftArithmetic2(func(x, y NNum) GNode { return x * y })
	case Div:
		s.liftArithmetic2(func(x, y NNum) GNode { return x / y })
	case Neg:
		s.push(s.Alloc(-s.Read(s.pop()).(NNum)))
	case Eq:
		s.liftCmp2(func(x, y NNum) bool { return x == y })
	case Ne:
		s.liftCmp2(func(x, y NNum) bool { return x != y })
	case Lt:
		s.liftCmp2(func(x, y NNum) bool { return x < y })
	case Le:
		s.liftCmp2(func(x, y NNum) bool { return x <= y })
	case Gt:
		s.liftCmp2(func(x, y NNum) bool { return x > y })
	case Ge:
		s.liftCmp2(func(x, y NNum) bool { return x >= y })

	default:
		panic("illegal instr")
	}
	return true
}

func (s *State) stepUnwind() {
	addr := s.pop()
	switch n := s.Read(addr).(type) {
	case NNum:
		// Eval 的结果是 Num, unwind num 代表 eval 结束
		if !s.Dump.empty() {
			dump := s.Dump.pop()
			// 对栈进行还原, 转回原代码执行
			s.Stack, s.Code = dump.Stack, dump.Code
		}
		s.push(addr)
	case NInd:
		// 将该间接节点内地址入栈，再次Unwind
		s.push(Addr(n))
		s.pushInstr(IUnwind)
	case NApp:
		// 将左侧地址入栈，再次Unwind
		s.push(addr)
		s.push(n.Fun)
		s.pushInstr(IUnwind)
	case NGlobal:
		// 在栈内有足够参数的情况下，将该超组合子加载到当前代码
		// 参数数量不足且 dump 中有保存的栈时，只保留原本的 redex 并且还原栈

		// weak head norm form
		partialApp := s.stackSz() < n.Arity
		if partialApp {
			// 保留redex, 还原栈
			// a1 : ...... : ak
			// ||
			// ak : s
			Assert(s.Dump.stackSz() > 0, "unwinding with too few args")
			dump := s.Dump.pop()
			s.drop(s.stackSz() - 1)
			s.pushN(dump.Stack.V...)
			s.Code = dump.Code
			return
		}

		if n.Arity == 0 {
			s.push(addr) // 留着 Global update
			s.pushInstr(n.Code...)
		} else {
			// 假设栈前面的 N 个地址指向一系列 NApp 节点
			// 保留最底部的一个(当作 Redex 更新用)
			// 清理掉上面N-1个地址，然后放上N个直接指向参数的地址
			apps := s.popN(n.Arity)
			args := make([]Addr, len(apps))
			for i, app := range apps {
				args[i] = s.Read(app).(NApp).Arg
			}

			// apps[0] 代表当前 unwind 的 Global app, e.g.
			// 	Stack
			//		App(App(Global(twice), Global(f)), Num(42))
			//		App(Global(twice), Global(f))
			//		Global(twice)
			//	Unwind
			// unwind 的 twice, 是 twice f 调用的 left most
			// ...
			// 	Stack
			//		App(App(Global(twice), Global(f)), Num(42))
			//		App(Global(twice), Global(f))
			//		Global(f)
			//		App(App(Global(compose), Global(f)), Global(f))
			//	Update(1), Pop(1), Unwind
			// compose f f  是 twice f 进行 unwind 的结果
			// 所以使用 compose f f 更新 twice f
			// 所以此处需要配合 compileSC, 把 unwind 的完整应用 apps[0] push 回去
			// push 的 args, 会被 compileSC 生成的 popN 清理掉
			s.push(apps[0])
			s.pushN(args...)
			s.pushInstr(n.Code...)
		}
	default:
		panic("unwind: wrong kind of node " + fmt.Sprint(n))
	}
}

func (s *State) liftArithmetic2(op func(lhs, rhs NNum) GNode) {
	lhs := s.Read(s.pop()).(NNum)
	rhs := s.Read(s.pop()).(NNum)
	s.push(s.Alloc(op(lhs, rhs)))
}
func (s *State) liftCmp2(op func(lhs, rhs NNum) bool) {
	lhs := s.Read(s.pop()).(NNum)
	rhs := s.Read(s.pop()).(NNum)
	if op(lhs, rhs) {
		s.push(s.Alloc(NNum(1)))
	} else {
		s.push(s.Alloc(NNum(0)))
	}
}

func (s *State) Reify() /*GNode*/ {
	fmt.Println(s.showCode())
	fmt.Println("-------------------")
	for s.step() {
		fmt.Println(s.showStack())
		fmt.Println(s.showCode())
		fmt.Println("-------------------")
	}
	//Assert(s.stackSz() == 1, "illegal state")
	//return s.Read(s.pop())
	Assert(s.stackSz() == 0, "illegal state") // print 消耗掉了?!
}

// 机器的初始状态下，所有编译好的超组合子都已经被放到堆上的NGlobal节点中，
// 而此时G-Machine中的当前代码序列只包含两条指令，
// 		第一条将main的对应节点地址放到栈上，
// 		第二条将main的对应指令序列加载到当前指令序列。
//
// main的对应指令序列会在堆上分配节点并装入相应数据，
// 最后在堆内存中构造出一个图，这个过程称为main的"实例化"。
// 构造完毕后这个图的入口地址会被放到栈顶。
// 完成实例化之后需要做收尾工作，即更新图节点(由于main没有参数，所以不必清理栈中的残留无用地址)并寻找下一个redex。

func Run(src string) []GNode {
	heap := NewHeap()
	globals := Globals{}

	for _, prim := range CompiledPrimitives {
		globals[prim.Name] = heap.Alloc(NGlobal(prim))
	}
	for _, sc := range Parse(PreludeDefs) {
		globals[sc.Name] = heap.Alloc(NGlobal(CompileSC(sc)))
	}
	for _, sc := range Parse(src) {
		globals[sc.Name] = heap.Alloc(NGlobal(CompileSC(sc)))
	}

	//var scs = append([]ScDef(nil), compiledPrimitives...)
	//for _, sc := range Parse(preludeDefs) {
	//	scs = append(scs, CompileSC(sc))
	//}
	//for _, sc := range Parse(src) {
	//	scs = append(scs, CompileSC(sc))
	//}
	//
	//for _, sc := range scs {
	//	globals[sc.Name] = heap.Alloc(NGlobal(sc))
	//}

	initCode := []Instr{
		PushGlobal("main"),
		IEval,
		IPrint,
	}

	state := State{
		Dump:    &Stack[Dump]{},
		Code:    initCode,
		Stack:   &Stack[Addr]{},
		Heap:    heap,
		Globals: globals,
	}
	state.Reify()
	return state.Output
}

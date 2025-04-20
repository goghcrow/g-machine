package lazy

import "fmt"

type Addr = int

type Node = any

const HNull = Addr(0)

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

// free 排序之后转换成 unused?

type Heap struct {
	free   []Addr
	unused Addr
	m      map[Addr]Node
}

func NewHeap() *Heap { return &Heap{unused: 1, m: map[Addr]Node{}} }

func (h *Heap) Alloc(n Node) (a Addr) {
	l := len(h.free)
	if l > 0 {
		a, h.free = h.free[l-1], h.free[:l-1]
	} else {
		a = h.unused
		h.unused++
	}
	h.m[a] = n
	return
}
func (h *Heap) Update(a Addr, n Node) {
	h.validate(a)
	h.m[a] = n
}
func (h *Heap) Free(a Addr) {
	h.validate(a)
	h.free = append(h.free, a)
	delete(h.m, a)
}
func (h *Heap) Lookup(a Addr) Node {
	n, ok := h.m[a]
	Assert(ok, "invalid addr")
	return n
}
func (h *Heap) Addresses() []Addr {
	xs := make([]Addr, 0, len(h.m))
	for a := range h.m {
		xs = append(xs, a)
	}
	return xs
}
func (h *Heap) Size() int { return len(h.m) }
func (h *Heap) validate(a Addr) {
	_, ok := h.m[a]
	Assert(ok, "invalid addr")
}

func IsNull(a Addr) bool {
	return a == HNull
}

func ShowAddr(a Addr) string {
	return fmt.Sprintf("#%d", a)
}

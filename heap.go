package g_machine

type Memory = []GNode

type Heap struct {
	Idx    int
	Memory Memory
}

func NewHeap() *Heap {
	return &Heap{
		Idx:    0,
		Memory: make([]GNode, 10000),
	}
}
func (h *Heap) Read(addr Addr) GNode        { return h.Memory[addr] }
func (h *Heap) Write(addr Addr, node GNode) { h.Memory[addr] = node }
func (h *Heap) Alloc(node GNode) Addr {
	for /*free*/ h.Memory[h.Idx] != nil {
		h.Idx = (h.Idx + 1) % len(h.Memory)
	}
	h.Memory[h.Idx] = node
	return h.Idx
}

// TODO GC
// https://amelia.how/posts/the-gmachine-in-detail.html#h7

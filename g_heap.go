package main

type Memory = []Node

type Heap struct {
	Idx    int
	Memory Memory
}

func (h *Heap) Read(addr Addr) Node        { return h.Memory[addr] }
func (h *Heap) Write(addr Addr, node Node) { h.Memory[addr] = node }
func (h *Heap) Alloc(node Node) Addr {
	for /*free*/ h.Memory[h.Idx] != nil {
		h.Idx = (h.Idx + 1) % len(h.Memory)
	}
	h.Memory[h.Idx] = node
	return h.Idx
}

// TODO GC
// https://amelia.how/posts/the-gmachine-in-detail.html#h7

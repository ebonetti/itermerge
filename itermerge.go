// Package itermerge is a package that provides primitives for an heap of iterators.
// It's a generic package when used in conjunction with github.com/taylorchu/generic
package itermerge

import (
	"container/heap"
)

// Type is the generic placeholder.
type Type looseLesser

// looseLesser is the generic placeholder interface, for Less method.
type looseLesser interface {
	Less(x interface{}) bool
}

// IterMerge is a min-heap of heap elements provided by user-defined iterator of sorted elements.
type IterMerge struct {
	h iterHeap
}

//IterMergeFrom creates a new heap of iterators from user-defined iterator of sorted elements.
func IterMergeFrom(nexts ...func() (info Type, ok bool)) *IterMerge {
	h := make(iterHeap, 0, len(nexts))
	for _, next := range nexts {
		if info, ok := next(); ok {
			h = append(h, iterator{info, next})
		}
	}
	heap.Init(&h)
	return &IterMerge{h}
}

// Push adds an iterator to an existing IterMerge.
func (m *IterMerge) Push(next func() (Type, bool)) {
	heap.Push(&m.h, next)
}

// Peek peeks the next element in the heap.
func (m IterMerge) Peek() (info Type, ok bool) {
	if len(m.h) > 0 {
		info, ok = m.h[0].Info, true
	}
	return
}

// Next returns the next element in the heap, if none ok is false; this method use the same logic as the iterators.
func (m *IterMerge) Next() (info Type, ok bool) {
	h := m.h

	if len(h) == 0 {
		return
	}

	info, ook := h[0].Next()
	h[0].Info, info = info, h[0].Info
	if ook {
		heap.Fix(&m.h, 0)
	} else {
		heap.Pop(&m.h)
	}

	return info, true
}

// iterHeap implemets code needed by container/heap.
type iterHeap []iterator

type iterator struct {
	Info Type
	Next func() (info Type, ok bool)
}

func (h iterHeap) Len() int {
	return len(h)
}

func (h iterHeap) Less(i, j int) bool {
	return h[i].Info.Less(h[j].Info)
}

func (h iterHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *iterHeap) Push(x interface{}) {
	next := x.(func() (Type, bool))
	if info, ok := next(); ok {
		*h = append(*h, iterator{info, next})
	}
}

func (h *iterHeap) Pop() interface{} {
	_h := *h
	n := len(_h)
	x := _h[n-1]
	_h[n-1] = iterator{}
	*h = _h[0 : n-1]
	return x
}

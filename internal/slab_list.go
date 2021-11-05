package internal

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Node struct {
	slab  []byte
	left  *Node
	right *Node
}

type SlabList struct {
	sem *semaphore.Weighted
	mux sync.Mutex
	ptr *Node
}

func NewSlabList(capacity int, slabSize int) *SlabList {
	sl := &SlabList{
		sem: semaphore.NewWeighted(int64(capacity)),
	}

	var prev *Node

	for i := 0; i < capacity; i++ {
		sl.ptr = &Node{make([]byte, slabSize), nil, prev}

		if prev != nil {
			prev.left = sl.ptr
		}

		prev = sl.ptr
	}

	sl.ptr = &Node{nil, nil, prev}
	prev.left = sl.ptr

	return sl
}

func (sl *SlabList) Push(slab []byte) {
	sl.mux.Lock()

	sl.ptr.slab = slab
	sl.ptr = sl.ptr.left

	sl.mux.Unlock()

	sl.sem.Release(1)
}

func (sl *SlabList) Pop() []byte {
	sl.sem.Acquire(context.Background(), 1)

	sl.mux.Lock()

	sl.ptr = sl.ptr.right
	slab := sl.ptr.slab

	sl.mux.Unlock()

	return slab
}

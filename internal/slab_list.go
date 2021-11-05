package internal

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Node struct {
	slab []byte
	next *Node
}

type SlabList struct {
	sem  *semaphore.Weighted
	head *Node
}

func NewSlabList(capacity int, slabSize int) *SlabList {
	sl := &SlabList{
		sem: semaphore.NewWeighted(int64(capacity)),
	}

	for i := 0; i < capacity; i++ {
		sl.head = &Node{make([]byte, slabSize), sl.head}
	}

	return sl
}

func (sl *SlabList) Push(slab []byte) {
	sl.sem.Acquire(context.Background(), 1)

	sl.head = &Node{slab, sl.head}

	sl.sem.Release(1)
}

func (sl *SlabList) Pop() []byte {
	sl.sem.Acquire(context.Background(), 1)

	slab := sl.head.slab
	sl.head = sl.head.next

	sl.sem.Release(1)

	return slab
}

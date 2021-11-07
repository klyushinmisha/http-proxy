package internal

import (
	"sync"
)

type SlabPool struct {
	pool sync.Pool
}

func NewSlabPool(capacity int, slabSize int) *SlabPool {
	sl := &SlabPool{}
	sl.pool.New = func() interface{} {
		slab := make([]byte, slabSize)
		return &slab
	}

	for i := 0; i < capacity; i++ {
		sl.pool.New()
	}

	return sl
}

func (sl *SlabPool) Put(slab *[]byte) {
	sl.pool.Put(slab)
}

func (sl *SlabPool) Get() *[]byte {
	return sl.pool.Get().(*[]byte)
}

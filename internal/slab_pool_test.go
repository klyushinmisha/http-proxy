package internal

import (
	"os"
	"sync"
	"testing"
)

func popPush(pool *SlabPool, wg *sync.WaitGroup) {
	pool.Put(pool.Get())
	wg.Done()
}

func Benchmark_SlabList_Alloc(b *testing.B) {
	const cap = 1000

	pool := NewSlabPool(os.Getpagesize())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}

		wg.Add(cap)

		for j := 0; j < cap; j++ {
			go popPush(pool, &wg)
		}

		wg.Wait()
	}
}

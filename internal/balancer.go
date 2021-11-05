package internal

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

type Balancer interface {
	Host(hosts []string) string
}

type roundRobinBalancer struct {
	ptr int64
}

func (b *roundRobinBalancer) Host(hosts []string) string {
	host := hosts[b.ptr]

	v := (int(b.ptr) + 1) % len(hosts)
	atomic.StoreInt64(&b.ptr, int64(v))

	return host
}

type randomBalancer struct {
	mux  sync.Mutex
	rand *rand.Rand
}

func (b *randomBalancer) Host(hosts []string) string {
	b.mux.Lock()

	i := b.rand.Int()

	b.mux.Unlock()

	return hosts[i%len(hosts)]
}

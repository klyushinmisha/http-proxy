package internal

import (
	"math/rand"
	"sync"
)

type Balancer interface {
	Host(hosts []string) string
}

type roundRobinBalancer struct {
	mux sync.Mutex
	pos int
}

func (b *roundRobinBalancer) Host(hosts []string) string {
	b.mux.Lock()

	i := b.pos
	b.pos = (b.pos + 1) % len(hosts)

	b.mux.Unlock()

	return hosts[i]
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

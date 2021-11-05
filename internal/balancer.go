package internal

import "math/rand"

type Balancer interface {
	Host(hosts []string) string
}

type roundRobinBalancer struct {
	ptr int
}

func (b *roundRobinBalancer) Host(hosts []string) string {
	host := hosts[b.ptr]
	b.ptr = (b.ptr + 1) % len(hosts)

	return host
}

type randomBalancer struct {
	rand *rand.Rand
}

func (b randomBalancer) Host(hosts []string) string {
	return hosts[rand.Int()%len(hosts)]
}

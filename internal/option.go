package internal

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	RoundRobinBalancer = "round-robin"
	RandomBalancer     = "random"
)

var (
	DefaultBalancer   Balancer = &roundRobinBalancer{}
	DefaultBufferSize          = os.Getpagesize()
)

func WithBalancerOfType(typ string) option {
	if typ == "" {
		return option{err: fmt.Errorf("balancer-type cannot be empty")}
	}

	switch typ {
	case RoundRobinBalancer:
		return option{
			balancer: &roundRobinBalancer{},
		}
	case RandomBalancer:
		return option{
			balancer: &randomBalancer{
				rand: rand.New(rand.NewSource(time.Now().Unix())),
			},
		}
	default:
		return option{err: fmt.Errorf("unknown balancer-type: %s", typ)}
	}
}

func WithBufferSize(s int) option {
	return option{bufferSize: s}
}

type option struct {
	balancer   Balancer
	bufferSize int
	err        error
}

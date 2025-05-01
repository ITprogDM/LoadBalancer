package balancer

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Balancer struct {
	counter  int
	backends []string
	mu       sync.Mutex
	log      *logrus.Logger
}

func NewBalancer(backends []string, log *logrus.Logger) *Balancer {
	return &Balancer{
		backends: backends,
		log:      log,
	}
}

func (b *Balancer) Next() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.backends) == 0 {
		b.log.Printf("Список серверов пуст")
		return ""
	}

	backend := b.backends[b.counter%len(b.backends)]
	b.counter++
	return backend
}

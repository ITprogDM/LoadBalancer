package balancer

import (
	"LoadBalancer/internal/balancer"
	"LoadBalancer/pkg/logger"
	"sync"
	"testing"
)

func TestRoundRobinBalancer(t *testing.T) {
	log := logger.NewLogger()
	backends := []string{"b1", "b2", "b3"}
	b := balancer.NewBalancer(backends, log)

	expected := []string{"b1", "b2", "b3", "b1", "b2"}
	for i, want := range expected {
		got := b.Next()
		if got != want {
			t.Errorf("RoundRobin step %d: got %s, want %s", i, got, want)
		}
	}
}

func BenchmarkBalancer_Parallel(b *testing.B) {
	log := logger.NewLogger()
	balancer := balancer.NewBalancer([]string{"b1", "b2", "b3", "b4"}, log)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = balancer.Next()
		}
	})
}

func TestBalancerRaceCondition(t *testing.T) {
	log := logger.NewLogger()
	balancer := balancer.NewBalancer([]string{"b1", "b2", "b3"}, log)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_ = balancer.Next()
			}
		}()
	}
	wg.Wait()
}

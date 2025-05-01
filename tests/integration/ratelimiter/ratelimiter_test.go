package ratelimiter

import (
	"LoadBalancer/internal/models"
	"LoadBalancer/internal/ratelimiter"
	"LoadBalancer/internal/repository"
	"context"
	"testing"
	"time"
)

type mockRepository struct {
	repo repository.ClientRepository
}

func (m *mockRepository) GetClient(_ context.Context, id string) (*models.Client, error) {
	return &models.Client{
		ID:         id,
		Capacity:   5,
		RatePerSec: 2,
	}, nil
}

// остальные методы не нужны для этого теста
func (m *mockRepository) UpsertClient(context.Context, *models.Client) error {
	return nil
}
func (m *mockRepository) DeleteClient(context.Context, string) error {
	return nil
}
func (m *mockRepository) ListClients(context.Context) ([]models.Client, error) {
	return nil, nil
}

func TestRateLimiterWithRepo(t *testing.T) {
	repo := &mockRepository{}
	rl := ratelimiter.NewRateLimiterWithRepository(repo)

	clientID := "test_client"

	// За короткий срок можем сделать не более 5 запросов (ёмкость)
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientID) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// 6-й запрос должен быть отклонён
	if rl.Allow(clientID) {
		t.Error("Expected rate limit to be exceeded")
	}

	// Ждём пополнение токенов
	time.Sleep(2 * time.Second)

	// Должен пройти
	if !rl.Allow(clientID) {
		t.Error("Expected request after refill to be allowed")
	}
}

func BenchmarkRateLimiterWithRepo(b *testing.B) {
	repo := &mockRepository{}
	rl := ratelimiter.NewRateLimiterWithRepository(repo)
	clientID := "bench_client"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rl.Allow(clientID)
		}
	})
}

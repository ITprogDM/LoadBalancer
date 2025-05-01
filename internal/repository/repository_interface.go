package repository

import (
	"LoadBalancer/internal/models"
	"context"
)

type ClientRepository interface {
	GetClient(ctx context.Context, id string) (*models.Client, error)
	UpsertClient(ctx context.Context, c *models.Client) error
	DeleteClient(ctx context.Context, id string) error
	ListClients(ctx context.Context) ([]models.Client, error)
}

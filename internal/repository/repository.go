package repository

import (
	"LoadBalancer/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewRepository(db *pgxpool.Pool, log *logrus.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func (s *Repository) GetClient(ctx context.Context, id string) (*models.Client, error) {
	row := s.db.QueryRow(ctx, "SELECT id, capacity, rate_per_sec FROM clients WHERE id=$1", id)
	var c models.Client
	err := row.Scan(&c.ID, &c.Capacity, &c.RatePerSec)
	return &c, err
}

func (s *Repository) UpsertClient(ctx context.Context, c *models.Client) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO clients (id, capacity, rate_per_sec)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) DO UPDATE SET
            capacity = EXCLUDED.capacity,
            rate_per_sec = EXCLUDED.rate_per_sec;
    `, c.ID, c.Capacity, c.RatePerSec)
	return err
}

func (s *Repository) DeleteClient(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, "DELETE FROM clients WHERE id=$1", id)
	return err
}

func (s *Repository) ListClients(ctx context.Context) ([]models.Client, error) {
	rows, err := s.db.Query(ctx, "SELECT id, capacity, rate_per_sec FROM clients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(&c.ID, &c.Capacity, &c.RatePerSec); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

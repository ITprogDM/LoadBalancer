package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func ClientPostgres(log *logrus.Logger) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, os.Getenv("PG_URL"))
	if err != nil {
		log.WithError(err).Error("Ошибка подключения к Postgres")
		return nil, err
	}

	log.Info("Успешное подключение к Postgres")
	return conn, nil
}

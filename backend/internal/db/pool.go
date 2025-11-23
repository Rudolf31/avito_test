package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewPool(lc fx.Lifecycle) *pgxpool.Pool {

	pool, errPool := pgxpool.New(context.Background(), "postgresql://postgres:postgres@postgres:5432/postgres?sslmode=disable")
	if errPool != nil {
		panic(errPool.Error())
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		panic(err.Error())
	}

	lc.Append(
		fx.StopHook(func() {
			pool.Close()
		}))

	return pool
}

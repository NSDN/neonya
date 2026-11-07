// Package postgres 只负责连接池的建立与关闭，不含任何业务逻辑。
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const pingTimeout = 5 * time.Second

// NewPool 建立连接池并确认数据库可连通。
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	poolConfiguration, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	poolConfiguration.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfiguration)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	pingContext, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := pool.Ping(pingContext); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

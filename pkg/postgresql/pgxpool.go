package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Client interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	Stat() *pgxpool.Stat
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// NewPostgresqlClient подключение psql-клиента.
func NewPostgresqlClient(ctx context.Context, url string, maxOpenConn, connMaxLifetimeMinute, connAttempts,
	connTimeout int) (*pgxpool.Pool, error) {

	connectConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	connectConfig.MaxConns = int32(maxOpenConn)
	connectConfig.MaxConnIdleTime = time.Minute * time.Duration(connMaxLifetimeMinute)

	var pool *pgxpool.Pool
	for connAttempts >= 0 {
		pool, err = pgxpool.NewWithConfig(ctx, connectConfig)
		if err == nil {
			break
		}

		time.Sleep(time.Second * time.Duration(connTimeout))

		connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("postgres connAttempts == 0: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

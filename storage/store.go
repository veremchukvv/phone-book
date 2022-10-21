package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type DB interface {
	ExecContext(ctx context.Context, sql string, args ...interface{}) (Result, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (Rows, error)
}

type Store struct {
	db      DB
	log     logrus.FieldLogger
	CloseFn func(ctx context.Context) error
}

func New(ctx context.Context, connConf *pgx.ConnConfig, log logrus.FieldLogger) (*Store, error) {
	conn, err := pgx.ConnectConfig(ctx, connConf)

	if err != nil {
		return nil, err
	}

	return &Store{
		db:  &PGConn{conn: conn},
		log: log,

		CloseFn: conn.Close,
	}, nil
}

type PGConn struct {
	conn *pgx.Conn
}

func (p *PGConn) ExecContext(ctx context.Context, sql string, args ...interface{}) (Result, error) {
	return p.conn.Exec(ctx, sql, args...)
}

func (p *PGConn) QueryContext(ctx context.Context, sql string, args ...any) (Rows, error) {
	return p.conn.Query(ctx, sql, args...)
}

type Row interface {
	Scan(...interface{}) error
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close()
}

type Result interface {
	RowsAffected() int64
}

func (s *Store) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return s.db.QueryContext(ctx, sql, args...)
}

func (s *Store) Exec(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	result, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

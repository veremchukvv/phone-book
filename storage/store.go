package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type DB interface {
	ExecContext(ctx context.Context, sql string, args ...interface{}) (Result, error)
	// ExecWithTransaction(ctx context.Context, sql string, txOptions pgx.TxOptions, args ...interface{}) (Result, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (Rows, error)
	StartTX(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	EndTX(ctx context.Context, tx pgx.Tx) error
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

// func (p *PGConn) ExecWithTransaction(ctx context.Context, sql string, txOptions pgx.TxOptions, args ...interface{}) (Result, error) {
// 	tx, err := p.conn.BeginTx(ctx, pgx.TxOptions{})

// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			log.Print(err)
// 			log.Print("Transaction Rollbacked")
// 			tx.Rollback(ctx)
// 		} else {
// 			tx.Commit(ctx)
// 			log.Print("Transaction Completed")
// 		}
// 	}()

// 	return tx.Exec(ctx, sql, args...)
// }

func (p *PGConn) StartTX(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	log.Print("Transaction opened")
	tx, err := p.conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		log.Print(err)
		return nil, err
	}
	// defer func() {
	// 	if err != nil {
	// 		log.Print(err)
	// 		log.Print("Transaction Rollbacked")
	// 		tx.Rollback(ctx)
	// 	} else {
	// 		tx.Commit(ctx)
	// 		log.Print("Transaction Completed")
	// 	}
	// }()

	return tx, nil
}

func (p *PGConn) EndTX(ctx context.Context, tx pgx.Tx) error {

	err := tx.Commit(ctx)
	if err != nil {
		return err
	}
	log.Print("Transaction Completed")

	return nil
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

func (s *Store) StartTX(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	tx, err := s.db.StartTX(ctx, pgx.TxOptions{})

	if err != nil {
		log.Print(err)
		return nil, err
	}
	// defer func() {
	// 	if err != nil {
	// 		log.Print(err)
	// 		log.Print("Transaction Rollbacked")
	// 		tx.Rollback(ctx)
	// 	} else {
	// 		tx.Commit(ctx)
	// 		log.Print("Transaction Completed")
	// 	}
	// }()

	return tx, nil
}

func (s *Store) EndTX(ctx context.Context, tx pgx.Tx) error {

	log.Print("Transaction ended")

	err := tx.Commit(ctx)
	if err != nil {
		return err
	}
	log.Print("Transaction Completed")

	return nil
}

func (s *Store) ExecContext(ctx context.Context, sql string, args ...interface{}) (Result, error) {
	return nil, nil
}

func (s *Store) QueryContext(ctx context.Context, sql string, args ...any) (Rows, error) {
	return nil, nil
}

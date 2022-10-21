package storage

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"regexp"
	"testing"
)

type DBMock struct {
	db *sql.DB
}

type RAffected struct {
	affected int64
}

func (r *RAffected) RowsAffected() int64 {
	return r.affected
}

type RowsMock struct {
	rows *sql.Rows
}

func (r *RowsMock) Scan(i ...interface{}) error {
	return r.rows.Scan(i...)
}

func (r *RowsMock) Next() bool {
	return r.rows.Next()
}

func (r *RowsMock) Close() {
	r.rows.Close()
}

func (d *DBMock) ExecContext(ctx context.Context, sql string, args ...interface{}) (Result, error) {
	return func() (Result, error) {
		res, err := d.db.ExecContext(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		count, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		return &RAffected{affected: count}, nil
	}()
}

func (d *DBMock) QueryContext(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return func() (Rows, error) {
		r, err := d.db.QueryContext(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		return &RowsMock{rows: r}, nil
	}()
}

func TestStore_FindUserByPhone(t *testing.T) {
	////
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Ошибка %s, создание мока", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "phone_number"}).
		AddRow("1", "Пользователь", "+793455555")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_id", "name", "phone_number" FROM "users" WHERE ("phone_number" = $1) LIMIT $2`)).
		WithArgs("+793455555", 1).
		WillReturnRows(rows)

	store := Store{
		db:  &DBMock{db: db},
		log: logrus.StandardLogger(),
		CloseFn: func(ctx context.Context) error {
			return db.Close()
		},
	}

	defer store.CloseFn(context.Background())
	////
	user, err := store.FindUserByPhone(context.Background(), "+793455555")
	if err != nil {
		t.Fatal(err)
	}

	if user != nil {
		return
	}

}

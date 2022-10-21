package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"integ/storage"
	"integ/storage/migration"
	"testing"
	"time"
)

type testCtx struct {
	t     *testing.T
	store *storage.Store
}

func prepareTestContext(t *testing.T) *testCtx {
	t.Helper()

	testCtx := &testCtx{t: t}

	store, cleanupFunc := prepareTestStore(t)
	t.Cleanup(cleanupFunc)

	testCtx.store = store

	return testCtx
}

func prepareTestStore(t *testing.T) (*storage.Store, func()) {
	t.Helper()

	connConfig, err := pgx.ParseConfig(GetPostgresURL())
	if err != nil {
		t.Fatal(err)
	}

	dbConn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	schemaName := fmt.Sprintf("test_%d", time.Now().UnixNano())
	t.Logf("Create %q DB schema for test", schemaName)
	if _, err = dbConn.Exec(context.Background(), "CREATE SCHEMA "+schemaName); err != nil {
		t.Fatalf("Failed to create test DB schema: %v", err)
	}
	connConfig.RuntimeParams["search_path"] = schemaName

	connString := fmt.Sprintf("%s search_path=%s", connConfig.ConnString(), schemaName)

	if err := evolution(connString); err != nil {
		t.Fatalf("ошибка при миграции бд %s", err)
	}

	store, err := storage.New(context.Background(), connConfig, logrus.StandardLogger())
	if err != nil {
		_ = dbConn.Close(context.Background())
		t.Fatalf("ошибка создания хранилища, причина %v", err)
	}

	return store, func() {
		defer dbConn.Close(context.Background())

		if _, err := dbConn.Exec(context.Background(), `DROP SCHEMA `+schemaName+` CASCADE`); err != nil {
			t.Errorf("drop schema %q failed: %v", schemaName, err)
		}
	}
}

func evolution(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	tooling := migration.New(db)

	return tooling.Run()
}

func TestMain(m *testing.M) {
	PostgresMain(m.Run)
}

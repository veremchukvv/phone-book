package migration

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"

	"github.com/pressly/goose"
)

type Tool struct {
	db *sql.DB
}

const (
	gooseTableName = "goose_migrations"
)

var (
	gooseMu = &sync.Mutex{}
)

func New(db *sql.DB, opts ...func(*Tool)) *Tool {
	tool := &Tool{
		db: db,
	}
	for _, opt := range opts {
		opt(tool)
	}
	return tool
}

func (t *Tool) Run() error {
	if err := t.runGooseMigration(); err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}
	return nil
}

func (t *Tool) runGooseMigration() error {
	gooseMu.Lock()
	defer gooseMu.Unlock()

	goose.SetTableName(gooseTableName)
	goose.SetVerbose(false)
	goose.SetLogger(logrus.New().WithField("subsys", "database_tool"))
	if err := goose.Up(t.db, "."); err != nil {
		return err
	}
	return nil
}

package migration

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInitialSchema, downInitialSchema)
}

func upInitialSchema(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE users 
(
	user_id SERIAL PRIMARY KEY,
	name VARCHAR(128) NOT NULL,
	phone_number TEXT NOT NULL
);

CREATE TABLE relations
(
    relation_id SERIAL PRIMARY KEY,
    user_id INTEGER not null,
    relation_user_id INTEGER not null 
);
`)
	return err
}

func downInitialSchema(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

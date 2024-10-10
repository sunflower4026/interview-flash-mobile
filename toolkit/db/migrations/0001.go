package migrations

import (
	"database/sql"
)

type V0001 struct{}

func init() {
	RegisterMigration(&V0001{})
}

func (v *V0001) Version() string {
	return "0001"
}

func (v *V0001) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		phone_number VARCHAR(255) NOT NULL UNIQUE,
		pin VARCHAR(255) NOT NULL,
		address VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );`)
	if err != nil {
		return err
	}
	return nil
}

func (v *V0001) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		return err
	}
	return nil
}

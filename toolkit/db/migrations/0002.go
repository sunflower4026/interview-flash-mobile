package migrations

import (
	"database/sql"
)

type V0002 struct{}

func init() {
	RegisterMigration(&V0002{})
}

func (v *V0002) Version() string {
	return "0002"
}

func (v *V0002) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
	ALTER TABLE users
	ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
	`)
	if err != nil {
		return err
	}
	return nil
}

func (v *V0002) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE users DROP COLUMN updated_at;`)
	if err != nil {
		return err
	}
	return nil
}

// migrations.go

// WARNING: DO NOT EDIT OR DELETE THIS FILE.
//
// This file contains sample of migration logic for the database and Elasticsearch.
// It is crucial for maintaining the integrity of the application's schema and data.

package migrations

import (
	"database/sql"
)

// V0000 is a sample migration struct
type V0000 struct{}

// Register your migration
func init() {
	RegisterMigration(&V0000{})
}

// Version returns the version of the migration
func (v *V0000) Version() string {
	return "0000"
}

// Up applies the migration
func (v *V0000) Up(tx *sql.Tx) error {
	// Add your migration logic here
	return nil
}

// Down reverts the migration
func (v *V0000) Down(tx *sql.Tx) error {
	// Add your rollback logic here
	return nil
}

package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

type Migration interface {
	Version() string
	Up(*sql.Tx) error
	Down(*sql.Tx) error
}

var migrationRegistry = make(map[string]Migration)

func RegisterMigration(m Migration) {
	migrationRegistry[m.Version()] = m
}
func LoadMigrations() ([]Migration, error) {
	var migrations []Migration
	for _, migration := range migrationRegistry {
		migrations = append(migrations, migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version() < migrations[j].Version()
	})

	return migrations, nil
}

func Migrate(postgresConf *sql.DB) error {
	// Load migrations outside of the switch statement
	migrations, err := LoadMigrations()
	if err != nil {
		return err
	}

	// Default to "up" migration if no arguments are provided
	if len(os.Args) < 2 {
		return MigrateUp(postgresConf, migrations)
	}

	// Process command line arguments
	switch os.Args[1] {
	case "up":
		return MigrateUp(postgresConf, migrations)
	case "down":
		return MigrateDown(postgresConf, migrations)
	case "rollback":
		if len(os.Args) < 3 {
			return fmt.Errorf("missing migration version for rollback")
		}
		return RollbackMigration(postgresConf, os.Args[2], migrations)
	default:
		return fmt.Errorf("invalid migration direction: %s", os.Args[1])
	}
}

func MigrateUp(db *sql.DB, migrations []Migration) error {
	if err := createVersionTable(db); err != nil {
		return fmt.Errorf("failed to create version table: %w", err)
	}

	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	latestVersion := migrations[len(migrations)-1].Version()

	log.Printf("Current applied version: %s\n", currentVersion)
	log.Printf("Latest available version: %s\n", latestVersion)

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, migration := range migrations {
		if !isMigrationApplied(tx, migration.Version()) {
			if err := migration.Up(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to apply migration %s: %w, rolling back", migration.Version(), err)
			}
			if err := recordMigration(tx, migration.Version()); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %w", migration.Version(), err)
			}
			log.Printf("Successfully applied migration: %s\n", migration.Version())
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Migrations up to version %s successfully applied\n", latestVersion)
	return nil
}

func MigrateDown(db *sql.DB, migrations []Migration) error {
	if err := createVersionTable(db); err != nil {
		return fmt.Errorf("failed to create version table: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		if isMigrationApplied(tx, migrations[i].Version()) {
			if err := migrations[i].Down(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to rollback migration %s: %w, rolling back", migrations[i].Version(), err)
			}
			if err := removeMigration(tx, migrations[i].Version()); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to remove migration %s: %w", migrations[i].Version(), err)
			}
			log.Printf("Successfully rolled back migration: %s\n", migrations[i].Version())
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Migrations down to version none successfully applied\n")

	return nil
}

func RollbackMigration(db *sql.DB, targetVersion string, migrations []Migration) error {
	if err := createVersionTable(db); err != nil {
		return fmt.Errorf("failed to create version table: %w", err)
	}

	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if currentVersion > targetVersion {
		// Rollback down to the target version
		for i := len(migrations) - 1; i >= 0; i-- {
			if migrations[i].Version() <= currentVersion && migrations[i].Version() > targetVersion {
				if isMigrationApplied(tx, migrations[i].Version()) {
					if err := migrations[i].Down(tx); err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to rollback migration %s: %w, rolling back", migrations[i].Version(), err)
					}
					if err := removeMigration(tx, migrations[i].Version()); err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to remove migration %s: %w", migrations[i].Version(), err)
					}
					log.Printf("Successfully rolled back migration: %s\n", migrations[i].Version())
				}
			}
		}
	} else {
		// Roll forward up to the target version
		for _, migration := range migrations {
			if migration.Version() > currentVersion && migration.Version() <= targetVersion {
				if !isMigrationApplied(tx, migration.Version()) {
					if err := migration.Up(tx); err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to apply migration %s: %w, rolling back", migration.Version(), err)
					}
					if err := recordMigration(tx, migration.Version()); err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to record migration %s: %w", migration.Version(), err)
					}
					log.Printf("Successfully applied migration: %s\n", migration.Version())
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func createVersionTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMPTZ NOT NULL
    );`)
	return err
}

func isMigrationApplied(tx *sql.Tx, version string) bool {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		log.Fatalf("failed to check migration version %s: %v", version, err)
	}
	return count > 0
}

func recordMigration(tx *sql.Tx, version string) error {
	_, err := tx.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)", version, time.Now())
	return err
}

func removeMigration(tx *sql.Tx, version string) error {
	_, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
	return err
}

func getCurrentVersion(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return "none", nil
	}
	if err != nil {
		return "", err
	}
	return version, nil
}

package goose

import (
	"database/sql"
	"sort"
)

func Reset(db *sql.DB, dir string) error {
	migrations, err := collectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}
	statuses, err := migrationsStatus(db)
	if err != nil {
		return err
	}
	sort.Sort(sort.Reverse(migrations))

	for _, migration := range migrations {
		if !statuses[migration.Version] {
			continue
		}
		if err = migration.Down(db); err != nil {
			return err
		}
	}

	return nil
}

package goose

import "database/sql"

func Apply(db *sql.DB, dir string) error {
	migrations, err := CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}
	statuses, err := migrationsStatus(db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if statuses[migration.Version] {
			continue
		}
		if err = migration.Up(db); err != nil {
			return err
		}
	}

	return nil
}

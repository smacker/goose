package migrations

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

var gDB *gorm.DB

func SetDB(db *gorm.DB) {
	gDB = db
}

func wrap(fn func(*gorm.DB) error) func(*sql.Tx) error {
	return func(_ *sql.Tx) error {
		tx := gDB.Begin()
		if err := fn(tx); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}
}

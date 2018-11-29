package dbman

import "github.com/jinzhu/gorm"

func (DBMan) Transaction(db *gorm.DB, f func(tx *gorm.DB) error) error {
	tx := db.Begin()
	err := f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return tx.Error
}

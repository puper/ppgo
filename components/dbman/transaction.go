package dbman

import "github.com/jinzhu/gorm"

func (DBMan) Transaction(db *gorm.DB, f func() error) error {
	tx := db.Begin()
	err := f()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

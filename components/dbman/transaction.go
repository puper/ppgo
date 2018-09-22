package dbman

import "github.com/jinzhu/gorm"

func (DBMan) Transaction(tx *gorm.DB, f func() error) error {
	err := f()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

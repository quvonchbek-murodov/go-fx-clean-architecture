package repositories

import "gorm.io/gorm"

type QueryOption func(*gorm.DB) *gorm.DB

func With(path string, conds ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(path, conds...)
	}
}

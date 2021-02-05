package otgorm

import (
	"github.com/go-gormigrate/gormigrate/v2"

	"gorm.io/gorm"
)

// MigrateFunc is the func signature for migrating.
type MigrateFunc func(*gorm.DB) error

// RollbackFunc is the func signature for rollbacking.
type RollbackFunc func(*gorm.DB) error

// Migration represents a database migration (a modification to be made on the database).
type Migration struct {
	// ID is the migration identifier. Usually a timestamp like "201601021504".
	ID string
	// Migrate is a function that will br executed while running this migration.
	Migrate MigrateFunc
	// Rollback will be executed on rollback. Can be nil.
	Rollback RollbackFunc
}

type Migrations struct {
	Db         *gorm.DB
	Collection []*Migration
}

func convert(old []*Migration) []*gormigrate.Migration {
	var out []*gormigrate.Migration
	for _, m := range old {
		out = append(out, &gormigrate.Migration{
			ID:       m.ID,
			Migrate:  gormigrate.MigrateFunc(m.Migrate),
			Rollback: gormigrate.RollbackFunc(m.Rollback),
		})
	}
	return out
}

func (m Migrations) Migrate() error {
	migration := gormigrate.New(m.Db, &gormigrate.Options{}, convert(m.Collection))
	return migration.Migrate()
}

func (m Migrations) Rollback(id string) error {
	migration := gormigrate.New(m.Db, &gormigrate.Options{}, convert(m.Collection))
	if id == "-1" {
		return migration.RollbackLast()
	}
	return migration.RollbackTo(id)
}

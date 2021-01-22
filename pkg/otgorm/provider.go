package otgorm

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/DoNewsCode/std/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseConf struct {
	DatabaseType    string
	Dsn         string
	TablePrefix string
}

func ProvideDialector(conf *DatabaseConf) (gorm.Dialector, error) {
	if conf.DatabaseType == "mysql" {
		return mysql.Open(conf.Dsn), nil
	}
	if conf.DatabaseType == "sqlite" {
		return sqlite.Open(conf.Dsn), nil
	}
	return nil, fmt.Errorf("unknow database type %s", conf.DatabaseType)
}

func ProvideGormConfig(l log.Logger, conf *DatabaseConf) *gorm.Config {
	return &gorm.Config{
		Logger:                                   &logging.GormLogAdapter{Logging: l},
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.TablePrefix,
		},
	}
}

func ProvideGormDB(dialector gorm.Dialector, config *gorm.Config, tracer opentracing.Tracer) (*gorm.DB, func(), error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, nil, err
	}
	AddGormCallbacks(db, tracer)
	return db, func() {
		if sqlDb, err := db.DB(); err == nil {
			sqlDb.Close()
		}
	}, nil
}
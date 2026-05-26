package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type Options struct {
	Driver     string
	PrimaryDB  string
	ReplicaDBs []string
}

func dialector(driver, dsn string) (gorm.Dialector, error) {
	switch driver {
	case "postgres", "postgresql":
		return postgres.Open(dsn), nil
	default:
		return nil, errors.New("db driver not supported")
	}
}

func Connect(opt *Options) (*gorm.DB, error) {
	primary, err := dialector(opt.Driver, opt.PrimaryDB)
	if err != nil {
		return nil, err
	}

	var replicas []gorm.Dialector
	for _, dsn := range opt.ReplicaDBs {
		d, err := dialector(opt.Driver, dsn)
		if err != nil {
			return nil, err
		}
		replicas = append(replicas, d)
	}

	conn, err := gorm.Open(primary, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open failed: %w", err)
	}

	if len(replicas) > 0 {
		err = conn.Use(
			dbresolver.Register(dbresolver.Config{
				Replicas: replicas,
				Policy:   dbresolver.RandomPolicy{},
			}).
				SetConnMaxIdleTime(30 * time.Minute).
				SetConnMaxLifetime(2 * time.Hour).
				SetMaxIdleConns(100).
				SetMaxOpenConns(500),
		)
		if err != nil {
			return nil, fmt.Errorf("dbresolver register failed: %w", err)
		}
	}

	return conn, nil
}

package dbconnections

import (
	"log"

	"golang-project-structure/config"
	"golang-project-structure/internal/models"
	"golang-project-structure/pkg/db"

	"gorm.io/gorm"
)

type AppDB struct {
	DB *gorm.DB
}

func NewAppDBConnection(cfg *config.Config) *AppDB {
	conn, err := db.Connect(&db.Options{
		Driver:     cfg.DB.Driver,
		PrimaryDB:  cfg.DB.DSN(),
		ReplicaDBs: cfg.DB.ReplicaDSNs(),
	})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := conn.AutoMigrate(&models.UserModel{}); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	return &AppDB{DB: conn}
}

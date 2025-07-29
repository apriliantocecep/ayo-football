package config

import (
	"database/sql"
	"github.com/apriliantocecep/ayo-football/services/team/internal/entity"
	"github.com/apriliantocecep/ayo-football/shared"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Database struct {
	DB   *gorm.DB
	Conn *sql.DB
}

func NewDatabase(vaultClient *shared.VaultClient) *Database {
	var gormDB *gorm.DB

	secret := utils.GetVaultSecretConfig(vaultClient)

	dsn := secret["DATABASE_URL"]
	if dsn == nil || dsn == "" {
		log.Fatalln("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn.(string)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(100)
	conn.SetConnMaxLifetime(time.Hour)

	// Auto Migrate
	err = db.AutoMigrate(
		&entity.Team{},
	)
	if err != nil {
		log.Fatalf("failed to migrate the entities: %v", err)
	}

	gormDB = db

	return &Database{
		DB:   gormDB,
		Conn: conn,
	}
}
